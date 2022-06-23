package storage

import (
	"fmt"
	"net"

	c "github.com/zuri03/GoCloudStore/constants"
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
		if err != nil {
			fmt.Printf("Error accepting connections => %s\n", err.Error())
			break
		}
		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	defer connection.Close()
	protocol := make([]byte, 3)

	connection.Read(protocol)

	fmt.Printf("Got Protocol => %s\n", string(protocol))

	switch string(protocol) {
	case c.GET_PROTOCOL:
		sendFileToClientHandler(connection)
	case c.SEND_PROTOCOL:
		storeFileHandler(connection)
	case c.DELETE_PROTOCOL:
	case c.ERROR_PROTOCOL:
	}

	return
}
