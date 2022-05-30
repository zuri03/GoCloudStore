package server

import (
	"fmt"
	"net"
	"os"
)

func sendFileToClient(fileName string, connection net.Conn) error {
	fmt.Printf("SENDING %s TO CLIENT\n", fileName)
	metaData, err := os.Stat(fileName)
	if err != nil {
		return err
	}

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	if metaData.Size() <= 1024 {
		buffer := make([]byte, metaData.Size())
		file.Read(buffer)
		connection.Write(buffer)
	}

	connection.Write([]byte("EOF"))
	return nil
}
