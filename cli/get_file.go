package cli

import (
	"bytes"
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
func getFileCommand(username string, password string, input []string, metaClient *MetadataServerClient) {
	key := input[0]
	record, err := metaClient.getFileRecord(username, password, key)
	if err != nil {
		fmt.Printf("Error getting file record: %s\n", err.Error())
		return
	}

	fmt.Printf("record => %+v\n", record)
	//connection, err := net.DialTCP("tcp", nil, dataNodeAddress)
	connection, err := net.DialTimeout("tcp", ":8000", time.Duration(10)*time.Second)
	defer connection.Close()

	connection.Write([]byte(c.GET_PROTOCOL))
	meta := FileMetaData{
		Username: username,
		FileName: record.MetaData.Name,
		Size:     record.MetaData.Size,
	}

	signal := make([]byte, 3)
	connection.Read(signal)
	fmt.Printf("Signal => %s\n", string(signal))

	if err := sendMetaDataToServer(meta, connection); err != nil {
		fmt.Printf("Error sending meta data: %s\n", err.Error())
		return
	}

	//Let the server know its okay to begin sending data
	connection.Write([]byte(c.PROCEED_PROTOCOL))

	if err := getFileDataFromServer(record.MetaData.Name, int(record.MetaData.Size),
		connection); err != nil {
		//Log error
		fmt.Printf("Error retreiving file data: %s\n", err.Error())
		return
	}

	fmt.Printf("SUCCESSFULLY RETREIVED FILE FROM SERVER")
}

//uses a connection object to retrieve byte data from a storage server and store it in a file
func getFileDataFromServer(fileName string, fileSize int, connection net.Conn) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return err
	}
	defer file.Close()

	if fileSize <= c.MAX_CACHE_BUFFER_SIZE {
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

	fileDataCacheBuffer := new(bytes.Buffer)
	readBuffer := make([]byte, c.TEMP_BUFFER_SIZE)
	writeBuffertoFile := func(buffer *bytes.Buffer, file *os.File) {
		file.Write(buffer.Bytes())
		buffer.Reset()
	}

	for {
		numOfBytes, err := connection.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				if numOfBytes > 0 {
					fileDataCacheBuffer.Write(readBuffer[:numOfBytes])
					writeBuffertoFile(fileDataCacheBuffer, file)
				}
				break
			}
			fmt.Println("ERROR OCCURRED RETURING")
			return err
		}

		if numOfBytes == 0 {
			fmt.Println("finished reading breaking")
			break
		}
		fmt.Printf("read => %d\n", numOfBytes)
		fmt.Printf("read => %s\n", string(readBuffer[:numOfBytes]))
		if _, err = fileDataCacheBuffer.Write(readBuffer[:numOfBytes]); err != nil {
			return err
		}

		if fileDataCacheBuffer.Len() > c.MAX_CACHE_BUFFER_SIZE {
			writeBuffertoFile(fileDataCacheBuffer, file)
		}
	}

	return nil
}
