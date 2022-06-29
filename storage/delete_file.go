package storage

import (
	"fmt"
	"net"
	"os"

	c "github.com/zuri03/GoCloudStore/common"
)

func deleteFileHandler(connection net.Conn) {
	meta, err := acceptFileMetaData(connection)
	if err != nil {
		fmt.Printf("Error accepting file meta data: %s\n", err.Error())
		connection.Write([]byte(c.ERROR_PROTOCOL))
		return
	}

	if err := deleteFileData(meta, connection); err != nil {
		fmt.Printf("Error deleting file data: %s\n", err.Error())
		connection.Write([]byte(c.ERROR_PROTOCOL))
		return
	}

	connection.Write([]byte(c.SUCCESS_PROTOCOL))
}

func deleteFileData(meta FileMetaData, connection net.Conn) error {
	directoryName := meta.Username
	filePath := fmt.Sprintf("%s/%s", directoryName, meta.FileName)

	if _, err := os.Stat(filePath); err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}
