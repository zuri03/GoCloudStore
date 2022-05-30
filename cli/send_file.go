package cli

import (
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
)

func sendFileCommand(input []string, connection net.Conn) error {
	file, metaData, err := getFileFromMemory(input)
	if err != nil {
		return err
	}

	//SND: short for send
	//Lets the server know the client is about to send a file
	connection.Write([]byte("SND"))
	if err := sendFileToServer(file, metaData, connection); err != nil {
		return err
	}
	return nil
}

func getFileFromMemory(input []string) (*os.File, fs.FileInfo, error) {
	fileName := input[0]
	fileExtension := filepath.Ext(fileName)
	if fileExtension != ".txt" && fileExtension != ".rtf" && fileExtension != ".pdf" {
		return nil, nil, fmt.Errorf("Not a text file")
	}

	fileMetaData, err := os.Stat(fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting file metadata => %s\n", err.Error())
	}

	if fileMetaData.IsDir() {
		return nil, nil, fmt.Errorf("Cannot send directory to server => %s\n", err.Error())
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("File %s does not exist", fileName)
	}

	return file, fileMetaData, nil
}
