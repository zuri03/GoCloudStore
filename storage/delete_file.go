package storage

import (
	"fmt"
	"net"

	c "github.com/zuri03/GoCloudStore/constants"
)

func deleteFileHandler(connection net.Conn) {
	meta, err := acceptFileMetaData(connection)
	if err != nil {
		fmt.Printf("Error accepting file meta data: %s\n", err.Error())
		return
	}

	if err := deleteFileData(meta, connection); err != nil {
		fmt.Printf("Error deleting file data: %s\n", err.Error())
		return
	}

	connection.Write([]byte(c.SUCCESS_PROTOCOL))
}

func deleteFileData(meta FileMetaData, connection net.Conn) error {
	return nil
}
