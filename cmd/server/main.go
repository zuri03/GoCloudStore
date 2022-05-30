package main

import (
	"fmt"
	"net"

	"github.com/zuri03/GoCloudStore/records"
	"github.com/zuri03/GoCloudStore/server"
)

// Implement secure connections
func main() {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Printf("Error starting listener => %s\n", err.Error())
		return
	}

	keeper := records.InitRecordKeeper()
	handler := server.Handler{Keeper: keeper}
	for {
		fmt.Println("Waiting for connections")
		connection, err := listener.Accept()

		if err != nil {
			fmt.Printf("Error accepting connections => %s\n", err.Error())
		}

		handler.HandleConnection(connection)
	}
}
