package main

import (
	"fmt"
	"net"

	"github.com/zuri03/GoCloudStore/server"
)

// Implement secure connections
func main() {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Printf("Error starting listener => %s\n", err.Error())
		return
	}

	for {
		fmt.Println("Waiting for connections")
		connection, err := listener.Accept()

		if err != nil {
			fmt.Printf("Error accepting connections => %s\n", err.Error())
		}

		server.HandleConnection(connection)
	}
}
