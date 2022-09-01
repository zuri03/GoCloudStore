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

	key := parseSendCommandInput(input[1:])

	//If no key is found then the default behavior is to set the key to filename
	if key == "" {
		key = fileInfo.Name()
	}

	err = metaClient.createFileRecord(owner, key, fileInfo.Name(), fileInfo.Size()) //For now just leave the key as the file name
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

//This command parses any flags and returns the value assigned to that flag
//Currently the only availabe flag is -k or key if there are anymore flags added in the future parse them with this function
func parseSendCommandInput(input []string) string {
	fmt.Printf("Parsing input\n")
	key := ""

	for idx, part := range input {
		fmt.Printf("%d => %s\n", idx, part)
		//If -k is found then the next element in the input array should be the key
		if part == "-k" {
			key = input[idx+1]
		}
	}
	fmt.Printf("Key is now: %s\n", key)
	return key
}

func getFileFromMemory(fileName string) (*os.File, fs.FileInfo, error) {
	fileExtension := filepath.Ext(fileName)
	if fileExtension != ".txt" && fileExtension != ".rtf" && fileExtension != ".pdf" {
		return nil, nil, fmt.Errorf("Not a text file")
	}

	fileMetaData, err := os.Stat(fileName)
	if err != nil {
		return nil, nil, err
	}

	if fileMetaData.IsDir() {
		return nil, nil, fmt.Errorf("Cannot send directory to server => %s\n", err.Error())
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}

	return file, fileMetaData, nil
}
