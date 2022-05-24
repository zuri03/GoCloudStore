package server

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

func checkUserCredentials(username string, password string) {
	//stub
}

func authenticateConnection(connectionScanner *bufio.Scanner, connection net.Conn) bool {
	fmt.Println("AUTHENTICATING CLIENT")
	connectionScanner.Scan()
	fmt.Println("GOT TOKEN")

	username := connectionScanner.Text()
	connectionScanner.Scan()

	password := connectionScanner.Text()

	//From here pass username and password to some authenticaiton service
	checkUserCredentials(username, password)
	connection.Write([]byte("OK"))
	return true
}

type FileMetaData struct {
	Size int64
	Name string
}

func acceptFileMetaData(connectionScanner *bufio.Scanner, connection net.Conn) (*FileMetaData, error) {
	gob.Register(new(FileMetaData))
	meta := &FileMetaData{}
	decoder := gob.NewDecoder(connection)
	fmt.Println("WAITING FOR GOB")
	err := decoder.Decode(meta)
	if err != nil {
		fmt.Printf("Error decoding gob => %s\n", err.Error())
		return nil, err
	}
	return meta, nil
}

func createFile(metaData *FileMetaData) (*os.File, error) {
	file, err := os.OpenFile(metaData.Name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return nil, err
	}
	fmt.Println("CREATED FILE")
	return file, nil
}

func acceptFileData(metaData *FileMetaData, file *os.File, connection net.Conn) error {
	if metaData.Size >= 1024 {
		//Read the data from the connection in a loop
	} else {
		dataBuffer := make([]byte, 1024)
		connection.Read(dataBuffer)
		fmt.Printf("Server got => %s\n", string(dataBuffer))
		result := bytes.TrimFunc(dataBuffer, func(r rune) bool {
			if r == 0 {
				return true
			}
			return false
		})
		fmt.Printf("result => %s\n", result)
		file.Write(dataBuffer)
	}
	return nil
}

func HandleConnection(connection net.Conn) {
	fmt.Println("Handling connection")
	defer connection.Close()

	connectionScanner := bufio.NewScanner(connection)

	authenticated := authenticateConnection(connectionScanner, connection)

	fmt.Println("exited auth")
	if !authenticated {
		connection.Write([]byte("User Not Found"))
		return
	}

	fmt.Println("About to accept meta data")

	metaData, err := acceptFileMetaData(connectionScanner, connection)
	if err != nil {
		fmt.Printf("Error while accepting meta data => %s\n", err.Error())
		return
	}

	file, err := createFile(metaData)
	if err != nil {
		fmt.Printf("Error creating file => %s\n", err.Error())
		return
	}

	fmt.Println("CREATED FILE")
	err = acceptFileData(metaData, file, connection)
}
