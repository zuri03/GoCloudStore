package server

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	//"os"
)

func authenticateConnection(connectionScanner *bufio.Scanner, connection net.Conn) bool {
	fmt.Println("AUTHENTICATING CLIENT")
	connectionScanner.Scan()
	fmt.Println("GOT TOKEN")

	username := connectionScanner.Text()
	fmt.Printf("Username => %s\n", username)
	connectionScanner.Scan()
	//connection.Write([]byte("OK"))

	password := connectionScanner.Text()
	fmt.Printf("Password => %s\n", password)

	//From here pass username and password to some authenticaiton service
	fmt.Println("SUCCESS ENDING SESSION")
	connection.Write([]byte("OK"))
	return true
}

func acceptFileMetaData(connectionScanner *bufio.Scanner, connection net.Conn) error {
	decoder := gob.NewDecoder(connection)
	var meta os.FileInfo
	fmt.Printf("Waiting for meta data \n")
	err := decoder.Decode(&meta)
	if err != nil {
		return fmt.Errorf("Error decoding file => %s\n", err.Error())
	}
	fmt.Printf("Got meta data => %d\n", meta.Size())
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

	//acceptFileMetaData(connectionScanner, connection)

	/*
		file, err := os.OpenFile("examle.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			fmt.Printf("Error occured while opening file => %s\n", err.Error())
			return
		}

		for connectionScanner.Scan() {
			text := connectionScanner.Text()

			fmt.Printf("Line => %s\n", text)

			if text == "EOF" {
				fmt.Println("End of file received")
				return
			}

			file.WriteString(text)
		}
	*/
}
