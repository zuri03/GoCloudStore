package cli

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	c "github.com/zuri03/GoCloudStore/common"
)

//Opens and closes connections and calls the file retreival functions
//in the correct order
//Handles any errors that may occur by notifying the user
func getFileCommand(owner string, input []string, metaClient *MetaDataClient) {
	key := input[0]
	record, err := metaClient.getFileRecord(owner, key)
	if err != nil {
		fmt.Printf("Error getting file record: %s\n", err.Error())
		return
	}
	//connection, err := net.DialTCP("tcp", nil, dataNodeAddress)
	connection, err := net.DialTimeout("tcp", ":8000", time.Duration(10)*time.Second)
	defer connection.Close()
	meta := c.FileMetaData{
		Owner: owner,
		Name:  record.Name,
		Size:  record.Size,
	}
	encoder, _ := newEncoderDecoder(connection)
	if err := sendMetaDataToServer(c.GET_FRAME, meta, encoder); err != nil {
		fmt.Printf("Error sending meta data: %s\n", err.Error())
		return
	}

	/*
		if _, err := acceptFrame(decoder, c.PROCEED_FRAME); err != nil {
			fmt.Printf("Error waiting for proceed: %s\n", err.Error())
			return
		}
	*/

	if err := sendFrame(c.PROCEED_FRAME, encoder); err != nil {
		fmt.Printf("Error sending proceed %s\n", err.Error())
		return
	}

	if err := getFileDataFromServer(record.Name, int(record.Size),
		connection, encoder); err != nil {
		//Log error
		fmt.Printf("Error retreiving file data: %s\n", err.Error())
		return
	}
	fmt.Println("Successfully retreived file from server")
}

//uses a connection to retrieve byte data from a storage server and store it in a file
func getFileDataFromServer(fileName string, fileSize int, connection net.Conn, encoder *gob.Encoder) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return err
	}
	defer file.Close()

	/*
		if fileSize <= c.MAX_BUFFER_SIZE {
			fmt.Println("Storing file in one chunk")
			buffer := make([]byte, fileSize)
			if _, err := connection.Read(buffer); err != nil {
				return err
			}
			if _, err := file.Write(buffer); err != nil {
				return err
			}
			return nil
		}
	*/
	//fileDataCacheBuffer := new(bytes.Buffer)
	/*
		writeBuffertoFile := func(buffer *bytes.Buffer, file *os.File) {
			file.Write(buffer.Bytes())
			buffer.Reset()
		}
	*/

	readBuffer := make([]byte, c.TEMP_BUFFER_SIZE)
	bufferedWriter := bufio.NewWriterSize(file, c.MAX_BUFFER_SIZE)
	defer bufferedWriter.Flush()
	fmt.Println("SENDING IN LOOP")
	for {
		numOfBytes, err := connection.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				if numOfBytes > 0 {
					if _, err := bufferedWriter.Write(readBuffer[:numOfBytes]); err != nil {
						return err
					}
				}
				break
			}
			fmt.Println("ERROR OCCURRED RETURING")
			return err
		}

		if _, err := bufferedWriter.Write(readBuffer[:numOfBytes]); err != nil {
			return err
		}

		fmt.Println("SENDING PROCEED FRAME")
		if err := sendFrame(c.PROCEED_FRAME, encoder); err != nil {
			return err
		}
	}

	return nil
}
