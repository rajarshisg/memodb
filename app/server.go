package main

import (
	"fmt"
	"net"
	"os"

	"redis-clone/app/internal/commands"
)

// handleConnection function handles an incoming client TCP connection requests by acknowledging it with a response.
func handleConnection(clientConn net.Conn) {
	// close the connection on function exit
	defer clientConn.Close()

	for {
		var buffer []byte = make([]byte, 10240)
		numBytes, err := clientConn.Read(buffer)
		if err != nil {
			fmt.Println("Error while reading data from client: ", err.Error())
			return
		}

		// if there are some incoming messages in buffer
		if numBytes > 0 {
			response, err := commands.HandleCommand(buffer[:numBytes])
			if err != nil {
				fmt.Println("Error while responding to client: ", err.Error())
				return
			}
			_, err = clientConn.Write([]byte(response))
			if err != nil {
				fmt.Println("Error while responding to client: ", err.Error())
				return
			}
		}
	}
	
}

func main() {
	tcpListener, err := net.Listen("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	} else {
		fmt.Println("Redis server is running on port 6379")
	}
	
	for {
		clientConn, err := tcpListener.Accept()
		if err != nil {
			fmt.Println("Error accepting clientConn: ", err.Error())
			os.Exit(1)
		}

		// each connection is handled in a separate thread
		go handleConnection(clientConn)
	}
	
}
