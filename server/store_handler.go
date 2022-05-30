package server

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"github.com/zuri03/GoCloudStore/records"
)

func storeFileFromClient(connectionScanner *bufio.Scanner, connection net.Conn, keeper records.RecordKeeper) error {
	fmt.Println("ACCEPTING META DATA")
	metaData, err := acceptFileMetaData(connectionScanner, connection)
	if err != nil {
		return err
	}

	fmt.Println("CREATE FILE")
	file, err := createFile(metaData)
	if err != nil {
		return err
	}

	//TODO: Use context to get owner and key?
	keeper.SetRecord("PLACEHOLDER", "PLACEHOLDER", metaData)

	fmt.Println("ACCEPT FILE DATA")
	err = acceptFileData(metaData, file, connection)
	if err != nil {
		return err
	}

	fmt.Println("SUCCESS")
	return nil
}

func acceptFileData(metaData *records.FileMetaData, file *os.File, connection net.Conn) error {
	if metaData.Size >= 1024 {
		//Read the data from the connection in a loop
	} else {
		dataBuffer := make([]byte, 1024)
		connection.Read(dataBuffer)
		result := trimBuffer(dataBuffer)
		file.Write(result)
	}
	return nil
}

func createFile(metaData *records.FileMetaData) (*os.File, error) {
	file, err := os.OpenFile(metaData.Name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755) //TODO: Determine correct permission
	if err != nil {
		return nil, err
	}
	return file, nil
}

func acceptFileMetaData(connectionScanner *bufio.Scanner, connection net.Conn) (*records.FileMetaData, error) {
	gob.Register(new(records.FileMetaData))
	meta := &records.FileMetaData{}
	decoder := gob.NewDecoder(connection)
	fmt.Println("WAITING FOR GOB")
	err := decoder.Decode(meta)
	if err != nil {
		return nil, err
	}
	return meta, nil
}
