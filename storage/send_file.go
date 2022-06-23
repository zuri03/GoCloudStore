package storage

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"

	c "github.com/zuri03/GoCloudStore/constants"
)

func sendFileToClientHandler(connection net.Conn) {
	meta, err := acceptFileMetaData(connection)
	if err != nil {
		fmt.Printf("Error accepting meta: %s\n", err.Error())
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
	/*

		fileReader := bufio.NewReader(file)
		connectionWriter := bufio.NewWriter(connection)
		//A strange way to hanlde any error in flush
		defer func() {
			if err := connectionWriter.Flush(); err != nil {
				e = err
			}
		}()
		fmt.Printf("writer size => %d\n", connectionWriter.Size())
		fmt.Printf("reader size => %d\n", fileReader.Size())

		if meta.Size <= int64(MAX_CACHE_BUFFER_SIZE) {
			fileData, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Printf("ERROR IN IOUTIL READ => %s\n", err.Error())
				return err
			}
			fmt.Printf("ioutil buffer => %s\n", string(fileData))
			fmt.Printf("SENDING FILE IN ONE CHUNCK")
			connectionWriter.Write(fileData)
			fmt.Printf("writer size => %d\n", connectionWriter.Size())
			fmt.Printf("writer buf size => %d\n", connectionWriter.Buffered())
			fmt.Printf("reader size => %d\n", fileReader.Size())
			return nil
		}
		return nil
	*/
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
