package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func sendFileCommand(username string, password string, input []string, client *MetadataServerClient) error {
	if len(input) != 2 {
		return fmt.Errorf("Incorrect amount of arguments for command \n Proper Send Command Usage: \n send ./example.txt key=foo")
	}

	//For now just ignore the actual file object
	//The file object will be used when we need to send the file data to the correct data server
	_, metadata, err := getFileFromMemory(input)
	if err != nil {
		return err
	}

	if !strings.Contains(input[1], "key=") {
		return fmt.Errorf("Missing key for record \n Proper Send Command Usage: \n send ./example.txt key=foo")
	}

	key := strings.Split(input[1], "=")[1]

	client.createFileRecord(username, password, key, metadata.Name(), metadata.Size())

	//SND: short for send
	//Lets the server know the client is about to send a file
	/*
		connection.Write([]byte("SND"))
		if err := sendFileToServer(file, metaData, connection); err != nil {
			return err
		}
	*/
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
