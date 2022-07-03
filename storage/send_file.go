package storage

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"

	c "github.com/zuri03/GoCloudStore/common"
)

func sendFileHandler(connection net.Conn, frame c.ProtocolFrame) {
	_, _ = newEncoderDecorder(connection)
	meta, err := decodeMetaData(frame)
	if err != nil {
		fmt.Printf("Error decoding meta: %s\n", err.Error())
		return
	}

	if err := sendFileDataToClient(meta, connection); err != nil {
		fmt.Printf("Error while reading file data: %s\n", err.Error())
		return
	}

	fmt.Println("Successfully sent file data")
}

func sendFileDataToClient(meta FileMetaData, connection net.Conn) error {
	file, err := openFile(meta.Username, meta.FileName)
	if err != nil {
		return err
	}
	defer file.Close()

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
