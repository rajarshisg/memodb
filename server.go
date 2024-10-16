package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"memodb/internal/commands"
	"memodb/internal/store"
)

// handleConnection function handles an incoming client TCP connection requests by acknowledging it with a response.
func handleConnection(clientConn net.Conn) {
	// close the connection on function exit
	defer clientConn.Close()
	
	for {
		var buffer []byte = make([]byte, 10240)
		numBytes, err := clientConn.Read(buffer)
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("Error reading data from client: ", err.Error())
			}
			break
		}

		// if there are some incoming messages in buffer
		if numBytes > 0 {
			response, err := commands.HandleCommand(buffer[:numBytes])
			if err != nil {
				fmt.Println("Error while responding to client: ", err.Error())
				break
			}
			_, err = clientConn.Write([]byte(response))
			if err != nil {
				fmt.Println("Error while responding to client: ", err.Error())
				break
			}
		}
	}
	
}

func main() {
	port := flag.String("port", "6379", "Port on which the Redis server runs")
	dir := flag.String("dir", "", "Path where RDB backups are stored")
	dbFileName := flag.String("dbfilename", "dump.rdb", "Name of the backup file")

	flag.Parse()

	commands.ConfigSet("dir", *dir)
	commands.ConfigSet("dbfilename", *dbFileName)
	_, err := store.LoadRdbInStore(*dir, *dbFileName)
	if (err != nil) {
		fmt.Printf("error occured while loadig rdb dump: %s\n", err.Error())
	}

	tcpListener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", *port))
	if err != nil {
		fmt.Printf("Failed to bind to port %s", *port)
		os.Exit(1)
	} else {
		fmt.Println()
		fmt.Println()
		fmt.Println("****************************************")
		fmt.Printf("* MemoDB server is running on port %s *\n", *port)
		fmt.Println("****************************************")
		fmt.Println()
		fmt.Println()
	}

	for {
		clientConn, err := tcpListener.Accept()
		if err != nil {
			fmt.Println("Error accepting clientConn: ", err.Error())
		}

		// each connection is handled in a separate thread
		go handleConnection(clientConn)
	}
	
}
