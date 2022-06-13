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
		defer func(connection net.Conn) {
			fmt.Printf("DONE READING SENDING SIGNAL")
			connection.Write([]byte("t"))
			fmt.Printf("SENT SIGNAL")
		}(connection)
			file, err := os.Open(fmt.Sprintf("./%s/%s", meta.Username, meta.FileName))
			if err != nil {
				return err
			}
	*/
	fileDataCacheBuffer := new(bytes.Buffer)
	readBuffer := make([]byte, meta.Size+100)
	fmt.Printf("initial length => %d\n", fileDataCacheBuffer.Len())
	if meta.Size <= int64(MAX_CACHE_BUFFER_SIZE) {
		fmt.Println("READING FILE IN ONE CHUNCK")
		/*
			num, err := fileDataCacheBuffer.ReadFrom(connection)
			if err != nil {
				fmt.Printf("Error: %d => %s\n", num, err.Error())
			}
			fmt.Printf("No error => %d => %d => %s\n", num, fileDataCacheBuffer.Len(), string(fileDataCacheBuffer.Bytes()))
		*/

		n, err := connection.Read(readBuffer)
		fmt.Printf("read buffer => %s\n", string(readBuffer))
		fmt.Printf("len => %d\n", n)
		if err != nil {
			fmt.Printf("read buffer error => %s\n", string(readBuffer))
			fmt.Printf("len error => %d\n", n)
			return err
		}

		return nil
	}
	fmt.Println("ABOUT TO BEGIN READING LOOP")

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
			/*
				file.Write(fileDataCacheBuffer.Bytes())
				fileDataCacheBuffer.Reset()
			*/
		}
	}
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
	fmt.Printf("Meta size => %s\n", meta.FileName)
	/*
		con := make([]byte, 1024)
		connection.Read(con)
		fmt.Printf("con => %s\n", string(con))
	*/
	return meta, nil
}
