package storage

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"

	c "github.com/zuri03/GoCloudStore/common"
)

func sendFileToClientHandler(connection net.Conn) {

	connection.Write([]byte(c.PROCEED_PROTOCOL))

	meta, err := acceptFileMetaData(connection)
	if err != nil {
		fmt.Printf("Error accepting meta: %s\n", err.Error())
		return
	}

	//Wait for signal from client to begin sending data
	signal := make([]byte, 3)
	connection.Read(signal)

	if string(signal) != c.PROCEED_PROTOCOL {
		fmt.Printf("Signal: %s\n", string(signal))
		fmt.Println("Error on client")
		return
	}

	if err := sendFileDataToClient(meta, connection); err != nil {
		fmt.Printf("Error while reading file data: %s\n", err.Error())
		return
	}

	fmt.Println("Successfully sent file data")
}

func sendFileDataToClient(meta FileMetaData, connection net.Conn) error {
	directoryName := meta.Username
	filePath := fmt.Sprintf("%s/%s", directoryName, meta.FileName)
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0400)
	if err != nil {
		return err
	}

	if meta.Size <= int64(MAX_CACHE_BUFFER_SIZE) {
		buffer := make([]byte, meta.Size)
		if _, err := file.Read(buffer); err != nil {
			return err
		}
		if _, err := connection.Write(buffer); err != nil {
			return err
		}
		return nil
	}

	bufferedConnWriter := bufio.NewWriter(connection)
	fileBuffer := make([]byte, c.TEMP_BUFFER_SIZE)
	for {
		numOfBytes, err := file.Read(fileBuffer)
		if err != nil {
			if err.Error() == "EOF" {
				if numOfBytes > 0 {
					bufferedConnWriter.Write(fileBuffer[:numOfBytes])
					if err := bufferedConnWriter.Flush(); err != nil {
						return nil
					}
				}
				return nil
			} else {
				return err
			}
		}

		bufferedConnWriter.Write(fileBuffer[:numOfBytes])

		if bufferedConnWriter.Buffered() >= c.MAX_CACHE_BUFFER_SIZE {
			if err := bufferedConnWriter.Flush(); err != nil {
				return err
			}
		}
	}
}

func acceptFileMetaData(connection net.Conn) (FileMetaData, error) {
	meta := FileMetaData{}
	decoder := gob.NewDecoder(connection)
	fmt.Println("WAITING FOR GOB")
	if err := decoder.Decode(&meta); err != nil {
		fmt.Printf("ERROR IN DECODER: %s\n", err.Error())
		return meta, err
	}

	return meta, nil
}
