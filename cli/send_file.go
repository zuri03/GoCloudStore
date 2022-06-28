package cli

import (
	"encoding/gob"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"time"

	c "github.com/zuri03/GoCloudStore/constants"
)

//Since fs.FileInfo cannot be encoded by
type FileMetaData struct {
	Size     int64
	FileName string
	Username string
}

func sendFileCommand(username string, password string, input []string, metaClient *MetadataServerClient) {
	fileName := input[0]
	file, fileInfo, err := getFileFromMemory(fileName)
	defer file.Close()
	if err != nil {
		fmt.Printf("Error getting file from memory: %s\n", err.Error())
		return
	}

	if fileInfo == nil || file == nil {
		fmt.Printf("Error could not find file %s\n", fileName)
		return
	}

	err = metaClient.createFileRecord(username, password, fileInfo.Name(), fileInfo.Name(), fileInfo.Size()) //For now just leave the key as the file name
	if err != nil {
		fmt.Printf("Error sending creating file record: %s\n", err.Error())
		return
	}

	//TODO: The address of the datanode must come from the record server
	//dataNodeAddress, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	//connection, err := net.DialTCP("tcp", nil, dataNodeAddress)
	connection, err := net.DialTimeout("tcp", ":8000", time.Duration(10)*time.Second)
	defer connection.Close()
	connection.Write([]byte(c.SEND_PROTOCOL))
	meta := FileMetaData{
		Username: username,
		FileName: fileInfo.Name(),
		Size:     fileInfo.Size(),
	}
	if err := sendMetaDataToServer(meta, connection); err != nil {
		fmt.Printf("Error sending meta data to server: %s\n", err.Error())
		return
	}

	signal := make([]byte, 3)
	connection.Read(signal)

	if string(signal) != c.PROCEED_PROTOCOL {
		fmt.Printf("Signal: %s\n", string(signal))
		fmt.Println("Error on server")
		return
	}

	if err := sendFileDataToServer(file, meta, connection); err != nil {
		fmt.Printf("Error sending file data to server: %s\n", err.Error())
		return
	}

	fmt.Println("Successfully sent file to server")
}

func sendFileDataToServer(file *os.File, meta FileMetaData, connection net.Conn) error {

	if meta.Size <= int64(c.MAX_CACHE_BUFFER_SIZE) {
		buffer := make([]byte, meta.Size)
		if _, err := file.Read(buffer); err != nil {
			return err
		}

		if _, err := connection.Write(buffer); err != nil {
			return err
		}
		return nil
	}

	buffer := make([]byte, c.MAX_CACHE_BUFFER_SIZE)

	for {
		numOfBytes, err := file.Read(buffer)

		if err != nil {
			if err.Error() == "EOF" {
				if numOfBytes > 0 {
					connection.Write(buffer[:numOfBytes])
				}
				return nil
			} else {
				return err
			}
		}

		connection.Write(buffer[:numOfBytes])
	}
}

func sendMetaDataToServer(meta FileMetaData, connection net.Conn) error {
	gob.Register(new(FileMetaData))
	encoder := gob.NewEncoder(connection)
	fmt.Println("Encoded gob")
	err := encoder.Encode(meta)
	if err != nil {
		return err
	}
	fmt.Println("SENT META DATA")
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
