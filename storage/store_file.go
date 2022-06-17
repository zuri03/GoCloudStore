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

	if err := storeFileDataFromClient(connection); err != nil {
		return err
	}
	return nil
}

func storeFileDataFromClient(connection net.Conn) error {
	meta := FileMetaData{}
	decoder := gob.NewDecoder(connection)
	fmt.Println("WAITING FOR GOB")
	if err := decoder.Decode(&meta); err != nil {
		fmt.Printf("ERROR IN DECODER: %s\n", err.Error())
		return err
	}

	filePath := fmt.Sprintf("/%s/%s", meta.Username, meta.FileName)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	/*
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			//file does not exist
		}
	*/

	fileDataCacheBuffer := new(bytes.Buffer)
	readBuffer := make([]byte, MAX_CACHE_BUFFER_SIZE) //setting the read buffer size to meta.Size+200 fixes the error for some reason
	if meta.Size <= int64(MAX_CACHE_BUFFER_SIZE) {
		fmt.Println("READING FILE IN ONE CHUNCK")
		n, err := connection.Read(readBuffer)
		if err != nil {
			fmt.Printf("read buffer error => %s\n", string(readBuffer))
			fmt.Printf("len error => %d\n", n)
			return err
		}
		fmt.Printf("BUFFER => %s\n", string(readBuffer))
		file.Write(readBuffer)
		return nil
	}

	fmt.Println("ABOUT TO BEGIN READING LOOP")

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
