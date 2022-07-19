package storage

import (
	"fmt"
	"io"
	"net"
	"os"

	c "github.com/zuri03/GoCloudStore/common"
)

func deleteFileHandler(connection net.Conn, frame c.ProtocolFrame) {

	fmt.Println("IN DELETE HANDLER")
	encoder, _ := newEncoderDecorder(connection)
	meta, err := decodeMetaData(frame)
	if err != nil {
		if err != io.EOF {
			if err := sendErrorFrame(encoder, err.Error()); err != nil {
				fmt.Printf("Error sending error frame: %s\n", err.Error())
			}
		}
		fmt.Printf("Error decoding meta: %s\n", err.Error())
		return
	}

	if err := deleteFileData(meta); err != nil {
		if err := sendErrorFrame(encoder, err.Error()); err != nil {
			fmt.Printf("Error sending error frame: %s\n", err.Error())
		}
		fmt.Printf("Error deleting file data: %s\n", err.Error())
		return
	}

	if err := sendSuccessFrame(encoder); err != nil {
		fmt.Printf("Error on success: %s\n", err.Error())
	}
}

func deleteFileData(meta c.FileMetaData) error {
	directoryName := meta.Owner
	filePath := fmt.Sprintf("%s/%s", directoryName, meta.Name)

	if _, err := os.Stat(filePath); err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}
