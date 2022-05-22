package server

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

func authenticateConnection(connectionScanner *bufio.Scanner, connection net.Conn) bool {
	fmt.Println("AUTHENTICATING CLIENT")
	connectionScanner.Scan()
	fmt.Println("GOT TOKEN")

	username := connectionScanner.Text()
	fmt.Printf("Username => %s\n", username)
	connectionScanner.Scan()

	password := connectionScanner.Text()
	fmt.Printf("Password => %s\n", password)

	//From here pass username and password to some authenticaiton service
	fmt.Println("SUCCESS ENDING SESSION")
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
	fmt.Printf("Result => %+v\n", meta)
	return meta, nil
}

func createFile(metaData *FileMetaData) (*os.File, error) {
	file, err := os.OpenFile(metaData.Name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		fmt.Printf("Error occured while opening file => %s\n", err.Error())
		return nil, err
	}
	return file, nil
}

func acceptFileData(metaData *FileMetaData, file *os.File, connection net.Conn) error {
	if metaData.Size >= 1024 {
		//Read the data from the connection in a loop
	} else {
		dataBuffer := make([]byte, 1024)
		connection.Read(dataBuffer)
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

	err = acceptFileData(metaData, file, connection)
}
