package cli

import (
	"encoding/gob"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/zuri03/GoCloudStore/common"
)

func sendFileCommand(owner string, input []string, metaClient *MetaDataClient) {
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

	encoder, decoder := newEncoderDecoder(connection)
	/*
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
	*/

	meta := common.FileMetaData{
		Owner: owner,
		Name:  fileInfo.Name(),
		Size:  fileInfo.Size(),
	}

	if err := sendMetaDataToServer(common.SEND_FRAME, meta, encoder); err != nil {
		fmt.Printf("Error sending meta data to server: %s\n", err.Error())
		return
	}

	if _, err := acceptFrame(decoder, common.PROCEED_FRAME); err != nil {
		fmt.Printf("Error waiting for proceed: %s\n", err.Error())
		return
	}

	if err := sendFileDataToServer(file, meta, connection, decoder); err != nil {
		fmt.Printf("Error sending file data to server: %s\n", err.Error())
		return
	}

	fmt.Println("Successfully sent file to server")
}

func sendFileDataToServer(file *os.File, meta common.FileMetaData, connection net.Conn, decoder *gob.Decoder) error {
	/*
		if meta.Size <= int64(common.MAX_BUFFER_SIZE) {
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
	*/
	fmt.Println("SENDING FILE DATA")
	buffer := make([]byte, common.TEMP_BUFFER_SIZE)
	for {
		numOfBytes, err := file.Read(buffer)

		if err != nil {
			if err == io.EOF {
				if numOfBytes > 0 {
					if _, err := connection.Write(buffer[:numOfBytes]); err != nil {
						return err
					}
				}
				return nil
			} else {
				return err
			}
		}

		if _, err := connection.Write(buffer[:numOfBytes]); err != nil {
			return err
		}
		fmt.Println("WAITING FOR PROCEED")
		if _, err := acceptFrame(decoder, common.PROCEED_FRAME); err != nil {
			return err
		}
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
