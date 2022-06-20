package storage

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func sendFileToClientHandler(connection net.Conn) error {
	meta, err := acceptFileMetaData(connection)
	if err != nil {
		return err
	}
	sendFileDataToClient(meta, connection)
	return nil
}

func sendFileDataToClient(meta FileMetaData, connection net.Conn) (e error) {
	file, err := os.OpenFile(meta.FileName, os.O_RDONLY, 0655)
	if err != nil {
		return err
	}
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
