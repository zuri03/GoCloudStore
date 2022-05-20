package cli

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io/fs"
	"net"
	"os"
)

func getServerResponse(scanner *bufio.Scanner) (string, error) {
	fmt.Printf("Waiting for message")
	scanner.Scan()
	message := scanner.Text()
	fmt.Printf("Connection successful message => %s\n", message)

	if message != "OK" {
		return "", fmt.Errorf("Error expected OK got => %s\n", message)
	}

	return message, nil
}

func authenticateSession(connection net.Conn, scanner *bufio.Scanner, username string, password string) bool {

	fmt.Println("Authenticating")

	connection.Write([]byte(username))

	connection.Write([]byte(password))

	_, err := getServerResponse(scanner)
	if err != nil {
		fmt.Printf("Error connecting: %s\n", err.Error())
		return false
	}

	fmt.Println("SUCCESSFUL AUTHENTICATION")

	return false
}

func sendMetaDataToServer(meta fs.FileInfo, connection net.Conn) error {

	fmt.Println("Gen gob")
	encoder := gob.NewEncoder(connection)

	err := encoder.Encode(meta)
	if err != nil {
		return fmt.Errorf("Error encoding meta data => %s\n", err.Error())
	}
	fmt.Println("Encoded meta data")
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
