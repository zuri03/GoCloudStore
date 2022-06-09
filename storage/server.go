package storage

import (
	"fmt"
	"net"
)

func InitializeListener() {
	address, err := net.ResolveTCPAddr("tcp", ":8000")
	if err != nil {
		fmt.Printf("Error in resolver: %s\n", err.Error())
		return
	}
	//In the future the port will come from environment vars/ command line args
	listener, err := net.ListenTCP("tcp", address)

	if err != nil {
		fmt.Printf("Error in listener: %s\n", err.Error())
		return
	}

	for {
		fmt.Println("Waiting for connections")
		connection, err := listener.AcceptTCP()
		fmt.Println("GOT CONNECTION")
		if err != nil {
			fmt.Printf("Error accepting connections => %s\n", err.Error())
			break
		}
		fmt.Println("PASSING TO HANDLER")
		go handleConnection(connection)
	}
}

func handleConnection(connection *net.TCPConn) {
	fmt.Println("Got protocol")
	protocol := make([]byte, 3)

	fmt.Println("About to read")

	connection.Read(protocol)
	defer connection.Close()

	fmt.Printf("Got Protocol => %s\n", string(protocol))
	return
}
