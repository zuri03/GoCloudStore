package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"os"
)

func authenticateSession(connection net.Conn, scanner *bufio.Scanner, username string, password string) bool {
	fmt.Println("Authenticating")
	connection.Write([]byte(username))
	connection.Write([]byte(password))

	buf := make([]byte, 2)
	connection.Read(buf)

	fmt.Printf("Got message from read buf => %s\n", string(buf))
	fmt.Println("SUCCESSFUL AUTHENTICATION")

	return false
}

type Test struct {
	Val string `json:"value"`
}

func sendMetaDataToServer(meta fs.FileInfo, connection net.Conn) error {

	fmt.Println("Generated gob")
	test := Test{Val: "example"}
	jsonBytes, err := json.Marshal(test)
	if err != nil {
		return err
	}

	connection.Write(jsonBytes)
	connection.Write([]byte("\n"))

	fmt.Println("Sent json")

	/*
		encoder := gob.NewEncoder(connection)
		fmt.Println("Connected gob to buffer")
		err := encoder.Encode(meta)
		if err != nil {
			fmt.Printf("Error in gob => %s\n", err.Error())
			return err
		}
		fmt.Println("Encoded meta data")
		fmt.Println("Encoded meta data")
	*/
	return nil
}

func sendDataToServer() {

}

func sendFileToServer(file *os.File, meta fs.FileInfo, connection net.Conn) error {

	fmt.Println("Sending meta data to server")
	err := sendMetaDataToServer(meta, connection)
	if err != nil {
		return err
	}

	return nil

	/*

		//Arbitrary buffer size
		//Determine best buffer size later
		buffer := make([]byte, 500)

		for {
			numOfBytes, err := file.Read(buffer)

			if err != nil {
				if err.Error() == "EOF" {
					fmt.Println("End of file found")
					connection.Write([]byte("EOF"))
					return nil
				} else {
					return fmt.Errorf("Error occured while reading file => %s\n", err.Error())
				}
			}

			fmt.Printf("Number of bytes => %d\n", numOfBytes)
			if numOfBytes == 0 {
				fmt.Println("Finished reading file")
				break
			}

			connection.Write(buffer)
		}
		return nil
	*/
}
