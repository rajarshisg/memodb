package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"memodb/internal/commands"
	"memodb/internal/store"
	"memodb/internal/tcp"
	"memodb/internal/worker"
)

// handleConnection function handles an incoming client TCP connection requests by acknowledging it with a response.
func handleConnection(clientConn net.Conn, persist bool, threadPool chan int) {
	if threadPool != nil {
		defer func() {
			<-threadPool
		}()
	}
	defer clientConn.Close()

	startTimestamp := time.Now().Unix() // timestamp at which we started dealing with the connection
	for {
		var buffer []byte = make([]byte, 10240)
		numBytes, err := clientConn.Read(buffer)
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("Error reading data from client: ", err.Error())
			}
		}

		// if there are some incoming messages in buffer
		if numBytes > 0 {
			startTimestamp = time.Now().Unix() // if client is sending some data we keep the connection alive
			isSuccess, propagateCommand, err := commands.HandleCommand(clientConn, buffer[:numBytes])
			if propagateCommand {
				worker.PropagateCommand(buffer)
			}
			if !isSuccess || err != nil {
				fmt.Println("Error while responding to client: ", err.Error())
			}
		}

		if (time.Now().Unix() - startTimestamp) >= 30 && !persist {
			// if we haven't receieved any data for more than 30 seconds and we do not want to persist we close the connection
			break;
		}
	}
	
}

func main() {
	port := flag.String("port", "6379", "Port on which the Redis server runs")
	dir := flag.String("dir", "", "Path where RDB backups are stored")
	dbFileName := flag.String("dbfilename", "", "Name of the backup file")
	replicaOf := flag.String("replicaof", "", "Host Port")

	flag.Parse()

	// Reading RDB Dump
	if (*dir != "") {
		commands.ConfigSet("dir", *dir)
	}
	if (*dbFileName != "") {
		commands.ConfigSet("dbfilename", *dbFileName)
		_, err := store.LoadRdbInStore(*dir, *dbFileName)
		if (err != nil) {
			fmt.Printf("error occured while loadig rdb dump: %s\n", err.Error())
		}
	}
	

	tcpListener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", *port))
	if err != nil {
		fmt.Printf("Failed to bind to port %s", *port, )
		os.Exit(1)
	} else {
		masterHost := ""
		masterPort := ""
		
		if *replicaOf != "" {
			masterHost = strings.Split(*replicaOf, " ")[0]
			masterPort = strings.Split(*replicaOf, " ")[1]
		}
		workerId, err := worker.InitWorker(*replicaOf != "", "127.0.0.1", *port, masterHost, masterPort)
		if (workerId == "" || err != nil) {
			fmt.Println("Could not connect to master. Exiting...")
			return
		}
		fmt.Println()
		fmt.Println()
		fmt.Println("****************************************")
		fmt.Printf("* MemoDB server is running on port %s *\n", *port)
		fmt.Println("****************************************")
		fmt.Println()
		fmt.Println()
	}

	// persistent slave connections
	for _, conn := range tcp.Connections {
		if !conn.Initialized {
			conn.Initialized = true;
			go handleConnection(conn.Connection, true, nil)
		}
			
	}

	threadPoolSize := 10
	threadPool := make(chan int, threadPoolSize)

	// new tcp connections
	for {
		clientConn, err := tcpListener.Accept()
		if err != nil {
			fmt.Println("Error accepting clientConn: ", err.Error())
		}

		// each connection is handled in a separate thread
		select {
			case threadPool <- 1: 
				go handleConnection(clientConn, false, threadPool)
			default:
        		fmt.Println("Connection rejected: Thread pool full")
        		clientConn.Close()
		}
	}
	
}
