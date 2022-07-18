package cli

import (
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/zuri03/GoCloudStore/common"
)

func sendFileCommand(owner string, input []string, metaClient *MetaDataClient) {
	fmt.Printf("owner => %s\n", owner)
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

	err = metaClient.createFileRecord(owner, fileInfo.Name(), fileInfo.Name(), fileInfo.Size()) //For now just leave the key as the file name
	if err != nil {
		fmt.Printf("Error sending creating file record: %s\n", err.Error())
		return
	}

	//TODO: The address of the datanode must come from the record server
	connection, err := net.DialTimeout("tcp", ":8000", time.Duration(10)*time.Second)
	defer connection.Close()

	encoder, decoder := newEncoderDecorder(connection)
	meta := common.FileMetaData{
		Owner: owner,
		Name:  fileInfo.Name(),
		Size:  fileInfo.Size(),
	}
	fmt.Printf("param meta => %+v\n", meta)
	if err := sendMetaDataToServer(common.SEND_FRAME, meta, encoder); err != nil {
		fmt.Printf("Error sending meta data to server: %s\n", err.Error())
		return
	}

	if err := waitForProceed(decoder); err != nil {
		fmt.Printf("Error waiting for proceed: %s\n", err.Error())
		return
	}

	if err := sendFileDataToServer(file, meta, connection); err != nil {
		fmt.Printf("Error sending file data to server: %s\n", err.Error())
		return
	}

	fmt.Println("Successfully sent file to server")
}

func sendFileDataToServer(file *os.File, meta common.FileMetaData, connection net.Conn) error {

	if meta.Size <= int64(common.MAX_CACHE_BUFFER_SIZE) {
		fmt.Println("Sending file in one chunck")
		buffer := make([]byte, meta.Size)
		if _, err := file.Read(buffer); err != nil {
			return err
		}

		if _, err := connection.Write(buffer); err != nil {
			return err
		}
		return nil
	}

	buffer := make([]byte, common.MAX_CACHE_BUFFER_SIZE)

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
