package cli

import (
	"bytes"
	"fmt"
	"net"
)

func getFileCommand(username string, password string, input []string, metaClient *MetadataServerClient) error {
	key := input[0]
	fmt.Printf("u len =. %d\n", len(username))
	fmt.Printf("p len =. %d\n", len(password))
	record, err := metaClient.getFileRecord(username, password, key)
	if err != nil {
		return err
	}
	fmt.Printf("record => %+v\n", record)
	return nil
}

func getFileDataFromServer(connection net.Conn) (*bytes.Buffer, error) {
	fileDataReader := make([]byte, 1024)
	fileDataBuffer := new(bytes.Buffer)
	for {
		_, err := connection.Read(fileDataReader)
		if err != nil {
			return nil, err
		}

		trimmed := bytes.TrimRightFunc(fileDataReader, func(r rune) bool {
			if r == 0 {
				return true
			}
			return false
		})

		if string(trimmed[len(trimmed)-3:]) == "EOF" {
			fileDataBuffer.Write(trimmed[:len(trimmed)-3])
			break
		}
		fileDataBuffer.Write(trimmed)
	}
	return fileDataBuffer, nil
}
