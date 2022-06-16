package storage

import (
	"fmt"
	"net"
)

const (
	GET_PROTOCOL    string = "GET"
	ERROR_PROTOCOL  string = "ERR"
	SEND_PROTOCOL   string = "SND"
	DELETE_PROTOCOL string = "DEL"
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

func handleConnection(connection net.Conn) {
	defer connection.Close()
	fmt.Println("Got protocol")
	protocol := make([]byte, 3)

	fmt.Println("About to read")

	connection.Read(protocol)

	fmt.Printf("Got Protocol => %s\n", string(protocol))

	switch string(protocol) {
	case GET_PROTOCOL:
	case SEND_PROTOCOL:
		err := storeFileHandler(connection)
		if err != nil {
			fmt.Printf("Error in handler: %s\n", err.Error())
			return
		}
		fmt.Println("SUCCESSFULLY STORED FILE")
	case DELETE_PROTOCOL:
	case ERROR_PROTOCOL:
	}

	return
}
