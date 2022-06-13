package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	//"os"
)

type FileMetaData struct {
	Size     int64
	FileName string
	Username string
}

const (
	TEMP_BUFFER_SIZE      int = 256  //TODO: Determine best temp buffer size
	MAX_CACHE_BUFFER_SIZE int = 1024 //TODO: Determine caching buffer size
)

func storeFileHandler(connection *net.TCPConn) error {
	meta, err := acceptFileMetaData(connection)
	if err != nil {
		return err
	}

	if err := storeFileData(meta, connection); err != nil {
		return err
	}
	return nil
}

func storeFileData(meta FileMetaData, connection *net.TCPConn) error {
	/*

		file, err := os.Open(fmt.Sprintf("./%s/%s", meta.Username, meta.FileName))
		if err != nil {
			return err
		}
	*/

	if meta.Size <= int64(MAX_CACHE_BUFFER_SIZE) {
		fmt.Println("READING FILE IN ONE CHUNCK")
		fileDataBuffer := make([]byte, meta.Size)
		if _, err := connection.Read(fileDataBuffer); err != nil {
			fmt.Println("error reading connection")
			return err
		}
		return nil
	}

	fileDataCacheBuffer := new(bytes.Buffer)
	readBuffer := make([]byte, TEMP_BUFFER_SIZE)
	fmt.Println("ABOUT TO BEGIN READING LOOP")
	for {
		numOfBytes, err := connection.Read(readBuffer)
		if err != nil {
			return err
		}
		if numOfBytes == 0 {
			break
		}
		_, err = fileDataCacheBuffer.Write(readBuffer[:numOfBytes])
		if err != nil {
			fmt.Println("Error on appending file buffer")
			return err
		}

		if fileDataCacheBuffer.Len() > MAX_CACHE_BUFFER_SIZE {
			/*
				file.Write(fileDataCacheBuffer.Bytes())
				fileDataCacheBuffer.Reset()
			*/
		}
	}
	fmt.Printf("Final result => %s\n", string(fileDataCacheBuffer.Bytes()))
	fmt.Printf("Final length => %d\n", fileDataCacheBuffer.Len())
	return nil
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
