package storage

import (
	"encoding/gob"
	"fmt"
	"net"

	c "github.com/zuri03/GoCloudStore/common"
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

	decoder := gob.NewDecoder(connection)
	encoder := gob.NewEncoder(connection)
	metaFrame, err := acceptFrame(decoder)
	if err != nil {
		if err := sendErrorFrame(encoder, "Error accepting meta data"); err != nil {
			//Log error
		}
		fmt.Printf("error: %s\n", err.Error())
		return
	}
	fmt.Println("Got Frame")
	fmt.Printf("Frame type => %d\n", metaFrame.Type)
	switch metaFrame.Type {
	case c.GET_FRAME:
	case c.SEND_FRAME:
		storeFileHandler(connection, metaFrame)
	case c.DELETE_FRAME:
	default:
		if err := sendErrorFrame(encoder, "Unrecognized action type"); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
	}
	return
}
