package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

type FileMetaData struct {
	Size     int64
	Name     string
	Username string
}

const (
	TEMP_BUFFER_SIZE      int = 256  //TODO: Determine best temp buffer size
	MAX_CACHE_BUFFER_SIZE int = 1024 //TODO: Determine caching buffer size
)

func storeFileHandler(connection *net.TCPConn) error {
	fmt.Println("IN STORAGE HANDLER")
	meta, err := acceptFileMetaData(connection)
	if err != nil {
		return err
	}
	fmt.Printf("GOT META DATA %s\n", meta.Name)
	return nil
}

func storeFile(connection *net.TCPConn) error {

	fmt.Println("ABOUT TO BEGIN READING")
	fileDataCacheBuffer := new(bytes.Buffer)
	readBuffer := make([]byte, TEMP_BUFFER_SIZE)

	for {
		numOfBytes, err := connection.Read(readBuffer)
		if err != nil {
			fmt.Println("ERROR OCCURRED RETURING")
			return err
		}

		_, err = fileDataCacheBuffer.Write(readBuffer[:numOfBytes])
		if err != nil {
			fmt.Println("Error on appending file buffer")
			return err
		}

		if fileDataCacheBuffer.Len() > MAX_CACHE_BUFFER_SIZE {

		}
	}
}

func acceptFileMetaData(connection net.Conn) (FileMetaData, error) {
	meta := FileMetaData{}
	decoder := gob.NewDecoder(connection)
	fmt.Println("WAITING FOR GOB")
	err := decoder.Decode(&meta)
	if err != nil {
		fmt.Printf("ERROR IN DECODER: %s\n", err.Error())
		return meta, err
	}

	return meta, nil
}

/*
func appendFileDatatoFile(fileData []byte) error {

}
*/
