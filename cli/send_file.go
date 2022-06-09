package cli

import (
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
)

func sendFileCommand(username string, password string, input []string, metaClient *MetadataServerClient) error {
	fileName := input[0]
	file, meta, err := getFileFromMemory(fileName)
	if err != nil {
		fmt.Println(err.Error())
	}
	if meta == nil || file == nil {
		return fmt.Errorf("Error could not find file %s\n", fileName)
	}
	err = metaClient.createFileRecord(username, password, meta.Name(), meta.Name(), meta.Size()) //For now just leave the key as the file name
	if err != nil {
		return err
	}

	//TODO: The address of the datanode must come from the record server
	//dataNodeAddress, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	//dataNodeAddress, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		return err
	}
	fmt.Println("DIALED CONNECTION")
	//connection, err := net.DialTCP("tcp", nil, dataNodeAddress)
	connection, err := net.Dial("tcp", ":8000")
	defer connection.Close()
	fmt.Println("ABOUT TO SEND")
	connection.Write([]byte(SEND_PROTOCOL))

	return nil
}

func getFileFromMemory(fileName string) (*os.File, fs.FileInfo, error) {
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
