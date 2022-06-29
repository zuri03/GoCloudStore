package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"

	c "github.com/zuri03/GoCloudStore/common"
)

//TODO: Implement connection reading with buffered io
type FileMetaData struct {
	Size     int64
	FileName string
	Username string
}

const (
	TEMP_BUFFER_SIZE      int = 256  //TODO: Determine best temp buffer size\
	MAX_CACHE_BUFFER_SIZE int = 1024 //TODO: Determine caching buffer size
)

func storeFileHandler(connection net.Conn) {

	meta, err := acceptFileMetaData(connection)
	if err != nil {
		fmt.Printf("Error accepting meta: %s\n", err.Error())
		return
	}
	if err := storeFileDataFromClient(meta, connection); err != nil {
		fmt.Printf("Error storing file data: %s\n", err.Error())
		return
	}

	fmt.Println("Successfully store file")
}

func storeFileDataFromClient(meta FileMetaData, connection net.Conn) error {

	directoryName := meta.Username
	filePath := fmt.Sprintf("%s/%s", directoryName, meta.FileName)
	if _, err := os.Stat(directoryName); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(directoryName, 0644); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return err
	}
	defer file.Close()

	/*
		In order to fix the EOF errors the client must wait before it begins to send file data
		this write to the connection serves as a signal from the server to the client letting the client
		know that the server is ready to begin receiving file data
	*/
	connection.Write([]byte(c.PROCEED_PROTOCOL))
	if meta.Size <= int64(MAX_CACHE_BUFFER_SIZE) {
		buffer := make([]byte, meta.Size)
		if _, err := connection.Read(buffer); err != nil {
			return err
		}
		if _, err := file.Write(buffer); err != nil {
			return err
		}
		return nil
	}

	fileDataCacheBuffer := new(bytes.Buffer)
	readBuffer := make([]byte, TEMP_BUFFER_SIZE)
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

		if _, err = fileDataCacheBuffer.Write(readBuffer[:numOfBytes]); err != nil {
			return err
		}

		if fileDataCacheBuffer.Len() > MAX_CACHE_BUFFER_SIZE {
			writeBuffertoFile(fileDataCacheBuffer, file)
		}
	}
	return nil
}

func acceptMetaData(connection net.Conn) (FileMetaData, error) {
	meta := FileMetaData{}
	decoder := gob.NewDecoder(connection)
	if err := decoder.Decode(&meta); err != nil {
		return meta, err
	}

	return meta, nil
}
