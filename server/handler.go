package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
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

type Test struct {
	Val string `json:"value"`
}

func acceptFileMetaData(connectionScanner *bufio.Scanner, connection net.Conn) error {
	var meta Test
	buffer := new(bytes.Buffer)
	b := make([]byte, 100)
	for {
		fmt.Println("Waiting for data")
		numOfBytes, err := connection.Read(b)
		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				return err
			}
		}

		if numOfBytes == 0 {
			break
		}

		if bytes.Contains(b, []byte("\n")) {
			break
		}
		buffer.Write(b)
	}
	result := bytes.TrimFunc(buffer.Bytes(), func(r rune) bool {
		fmt.Printf("%c => %t\n", r, r == 0)
		if r == 0 {
			return true
		}
		return false
	})
	json.Unmarshal(result, &meta)
	fmt.Println("Got meta")
	fmt.Printf("meta => %+v\n", meta)
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

	//connection.Write([]byte("OK after meta"))

	acceptFileMetaData(connectionScanner, connection)

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
