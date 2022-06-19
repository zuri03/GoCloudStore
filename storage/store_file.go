package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
)

//TODO: Implement connection reading with buffered io
type FileMetaData struct {
	Size     int64
	FileName string
	Username string
}

const (
	TEMP_BUFFER_SIZE      int = 256  //TODO: Determine best temp buffer size
	MAX_CACHE_BUFFER_SIZE int = 1024 //TODO: Determine caching buffer size
)

func storeFileHandler(connection net.Conn) error {

	meta, err := acceptFileMetaData(connection)
	if err != nil {
		return err
	}
	if err := storeFileDataFromClient(meta, connection); err != nil {
		return err
	}
	return nil
}

func storeFileDataFromClient(meta FileMetaData, connection net.Conn) error {

	directoryName := meta.Username
	filePath := fmt.Sprintf("%s/%s", directoryName, meta.FileName)

	if _, err := os.Stat(directoryName); err != nil {
		if os.IsNotExist(err) {
			//file does not exist
			fmt.Printf("Directory does not exist %s\n", directoryName)
			if err := os.Mkdir(directoryName, 0644); err != nil {
				fmt.Printf("Error creating directory: %s\n", directoryName)
				return err
			}
			fmt.Printf("created directory for new user")
		} else {
			return err
		}
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return err
	}

	/*
		In order to fix the EOF errors the client must wait before it begins to send file data
		this write to the connection serves as a signal from the server to the client letting the client
		know that the server is ready to begin receiving file data
	*/
	connection.Write([]byte(PROCEED_PROTOCOL))

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

	fmt.Println("ABOUT TO BEGIN READING LOOP")
	fileDataCacheBuffer := new(bytes.Buffer)
	readBuffer := make([]byte, TEMP_BUFFER_SIZE)
	_ = func(buffer *bytes.Buffer, file *os.File) {
		file.Write(fileDataCacheBuffer.Bytes())
		fileDataCacheBuffer.Reset()
	}

	for {

		numOfBytes, err := connection.Read(readBuffer)
		if err != nil {
			if numOfBytes > 0 {
				fileDataCacheBuffer.Write(readBuffer[:numOfBytes])
				//resetCacheBuffer(fileDataCacheBuffer, file)
			}
			if err == io.EOF {
				break
			}
			fmt.Println("ERROR OCCURRED RETURING")
			return err
		}

		if numOfBytes == 0 {
			break
		}

		if _, err = fileDataCacheBuffer.Write(readBuffer[:numOfBytes]); err != nil {
			fmt.Println("Error on appending file buffer")
			return err
		}

		if fileDataCacheBuffer.Len() > MAX_CACHE_BUFFER_SIZE {
			//resetCacheBuffer(fileDataCacheBuffer, file)
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
