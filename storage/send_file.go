package storage

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"

	c "github.com/zuri03/GoCloudStore/common"
)

func sendFileHandler(connection net.Conn, frame c.ProtocolFrame) {
	encoder, _ := newEncoderDecorder(connection)
	meta, err := decodeMetaData(frame)
	if err != nil {
		fmt.Printf("Error decoding meta: %s\n", err.Error())
		if err := sendErrorFrame(encoder, err.Error()); err != nil {
			fmt.Printf("Error sending error frame: %s\n", err.Error())
		}
		return
	}

	file, err := openFileForTransfer(meta)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err.Error())
		if err := sendErrorFrame(encoder, err.Error()); err != nil {
			fmt.Printf("Error sending error frame: %s\n", err.Error())
		}
		return
	}
	defer file.Close()

	if err := sendFileDataToClient(file, meta, connection); err != nil {
		fmt.Printf("Error while reading file data: %s\n", err.Error())
		if err := sendErrorFrame(encoder, err.Error()); err != nil {
			fmt.Printf("Error sending error frame: %s\n", err.Error())
		}
		return
	}

	fmt.Println("Successfully sent file data")
}

func openFileForTransfer(meta c.FileMetaData) (*os.File, error) {
	directoryName := meta.Owner
	filePath := fmt.Sprintf("%s/%s", directoryName, meta.Name)

	var file *os.File
	if _, err := os.Stat(directoryName); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(directoryName, 0777); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func sendFileDataToClient(file *os.File, meta c.FileMetaData, connection net.Conn) error {

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

func acceptFileMetaData(connection net.Conn) (c.FileMetaData, error) {
	meta := c.FileMetaData{}
	decoder := gob.NewDecoder(connection)
	fmt.Println("WAITING FOR GOB")
	if err := decoder.Decode(&meta); err != nil {
		fmt.Printf("ERROR IN DECODER: %s\n", err.Error())
		return meta, err
	}

	return meta, nil
}
