package cli

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

func getFileCommand(fileName string, connection net.Conn) error {
	connection.Write([]byte(fmt.Sprintf("GET:%s", fileName)))
	fileData, err := getFileDataFromServer(connection)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0755)
	buf := make([]byte, fileData.Len())
	_, err = fileData.Read(buf)
	if err != nil {
		return err
	}
	file.Write(buf)
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
