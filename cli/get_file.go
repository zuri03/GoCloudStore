package cli

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

func getFileCommand(input []string, connection net.Conn) error {
	fileName := input[0]
	fmt.Printf("RETREIVING %s FROM SERVER \n", fileName)
	responseBuffer := make([]byte, 2)
	connection.Write([]byte(fmt.Sprintf("GET:%s", fileName)))
	fmt.Println("SENT PROTOCOL WAITING FOR RESPONSE")
	connection.Read(responseBuffer)
	fmt.Printf("GOT RESPONSE %s\n", string(responseBuffer))
	if string(responseBuffer) != "OK" {
		return fmt.Errorf("File %s not found\n", fileName)
	}

	fileData, err := getFileDataFromServer(connection)
	if err != nil {
		return err
	}

	fmt.Println("SAVING BUFFER TO FILE")
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0755)
	file.Write(fileData.Bytes())
	return nil
}

func getFileDataFromServer(connection net.Conn) (*bytes.Buffer, error) {
	fileDataReader := make([]byte, 1024)
	fileDataBuffer := new(bytes.Buffer)
	fmt.Println("READING FILE DATA FROM SERVER")
	for {
		connection.Read(fileDataReader)

		//If the last byte is empty just assume we have recieved all of the file contents
		if fileDataReader[len(fileDataReader)-1] == 0 {
			fileDataBuffer.Write(bytes.TrimRightFunc(fileDataReader, func(r rune) bool {
				if r == 0 {
					return true
				}
				return false
			}))
			break
		}

		fmt.Println("STORING BYTES IN BUFFER")
		fileDataBuffer.Write(fileDataReader)
	}

	return fileDataBuffer, nil
}
