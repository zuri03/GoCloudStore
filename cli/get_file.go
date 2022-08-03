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

	connection, err := net.DialTimeout("tcp", ":8000", time.Duration(10)*time.Second)
	defer connection.Close()
	meta := c.FileMetaData{
		Owner: owner,
		Name:  record.Name,
		Size:  record.Size,
	}
	encoder, decoder := newEncoderDecoder(connection)
	if err := sendMetaDataToServer(c.GET_FRAME, meta, encoder); err != nil {
		fmt.Printf("Error sending meta data: %s\n", err.Error())
		return
	}

	fmt.Println("Waiting for proceed")

	if _, err := acceptFrame(decoder, c.PROCEED_FRAME); err != nil {
		fmt.Printf("Error waiting for proceed: %s\n", err.Error())
		return
	}

	if err := getFileDataFromServer(meta.Name, int(meta.Size), connection, encoder); err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	if err := sendFrame(c.SUCCESS_FRAME, encoder); err != nil {
		fmt.Printf("Error on success: %s\n", err.Error())
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

	if fileSize <= c.MAX_BUFFER_SIZE {
		fmt.Println("Storing file in one chunk")
		buffer := make([]byte, fileSize)
		numOfBytes, err := connection.Read(buffer)
		if err != nil {
			return err
		}
		if _, err := file.Write(buffer[:numOfBytes]); err != nil {
			return err
		}
		return nil
	}

	readBuffer := make([]byte, c.TEMP_BUFFER_SIZE)
	bufferedWriter := bufio.NewWriterSize(file, c.MAX_BUFFER_SIZE)
	defer bufferedWriter.Flush()
	fmt.Println("SENDING IN LOOP")
	totalBytes := 0
	for totalBytes <= fileSize {
		fmt.Println("READY TO READ")
		numOfBytes, err := connection.Read(readBuffer)
		if err != nil && err != io.EOF {
			return err
		}

		if numOfBytes == 0 || err == io.EOF {
			fmt.Printf("GOT ALL OF THE BYTES BREAKING \n")
			//Just in case there are some bytes left over
			if numOfBytes != 0 {
				if _, err := bufferedWriter.Write(readBuffer[:numOfBytes]); err != nil {
					fmt.Printf("Error writing to final file: %s\n", err.Error())
					return err
				}
			}

			break
		}

		fmt.Printf("read %d bytes from connection: %s\n", numOfBytes, string(readBuffer[:numOfBytes]))
		if _, err := bufferedWriter.Write(readBuffer[:numOfBytes]); err != nil {
			fmt.Printf("error writing to file: %s\n", err.Error())
			return err
		}

		totalBytes += numOfBytes
	}

	return nil
}
