package cli

import (
	"encoding/gob"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
)

const (
	BLOCK_SIZE int = 1024
)

//Since fs.FileInfo cannot be encoded by
type FileMetaData struct {
	Size     int64
	FileName string
	Username string
}

func sendFileCommand(username string, password string, input []string, metaClient *MetadataServerClient) error {
	fileName := input[0]
	file, fileInfo, err := getFileFromMemory(fileName)
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	if fileInfo == nil || file == nil {
		return fmt.Errorf("Error could not find file %s\n", fileName)
	}
	err = metaClient.createFileRecord(username, password, fileInfo.Name(), fileInfo.Name(), fileInfo.Size()) //For now just leave the key as the file name
	if err != nil {
		return err
	}

	//TODO: The address of the datanode must come from the record server
	//dataNodeAddress, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	//dataNodeAddress, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		return err
	}
	//connection, err := net.DialTCP("tcp", nil, dataNodeAddress)
	connection, err := net.Dial("tcp", ":8000")
	defer connection.Close()
	connection.Write([]byte(SEND_PROTOCOL))
	meta := FileMetaData{
		Username: username,
		FileName: fileInfo.Name(),
		Size:     fileInfo.Size(),
	}

	if err := sendMetaDataToServer(meta, connection); err != nil {
		return nil
	}

	if err := sendFileDataToServer(file, meta, connection); err != nil {
		return err
	}
	return nil
}

func sendFileDataToServer(file *os.File, meta FileMetaData, connection net.Conn) error {
	/*
		defer func(connection net.Conn) {
			fmt.Println("WAITING FOR GO AHEAD TO RETURN")
			signal := make([]byte, 1)
			connection.Read(signal)
			fmt.Println("SIGNAL RECEIVED CLOSING CONNECTION")
		}(connection)
	*/

	if meta.Size <= int64(BLOCK_SIZE) {
		fmt.Printf("SENDING FILE IN ONE CHUNK => %d\n", meta.Size)
		buffer := make([]byte, meta.Size)
		n, e := file.Read(buffer)
		if e != nil {
			fmt.Printf("error reading file => %s\n", e.Error())
			return e
		}
		fmt.Printf("file => %d\n", n)
		fmt.Printf("buffer => %s\n", string(buffer))
		n, err := connection.Write(buffer)
		if err != nil {
			return err
		}
		fmt.Printf("SENT FILE DATA => %d\n", n)

		return nil
	}

	buffer := make([]byte, BLOCK_SIZE)

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
}

func sendMetaDataToServer(meta FileMetaData, connection net.Conn) error {
	fmt.Println("Generated gob")
	gob.Register(new(FileMetaData))
	encoder := gob.NewEncoder(connection)
	fmt.Println("Connected gob to buffer")
	err := encoder.Encode(meta)
	if err != nil {
		fmt.Printf("ERORR ENCODING: %s\n", err.Error())
		return err
	}
	fmt.Println("Encoded meta data")
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
