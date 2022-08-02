package storage

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"

	c "github.com/zuri03/GoCloudStore/common"
)

func sendFileHandler(connection net.Conn, frame c.ProtocolFrame) {
	encoder, decoder := newEncoderDecorder(connection)
	meta, err := decodeMetaData(frame)
	fmt.Printf("Meta Data: %+v\n", meta)
	if err != nil {
		fmt.Printf("Error decoding meta: %s\n", err.Error())
		if err := sendErrorFrame(encoder, err.Error()); err != nil {
			fmt.Printf("Error sending error frame: %s\n", err.Error())
		}
		return
	}

	file, err := openFileForTransfer(meta)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err.Error())
		if err := sendErrorFrame(encoder, err.Error()); err != nil {
			fmt.Printf("Error sending error frame: %s\n", err.Error())
		}
		return
	}
	defer file.Close()

	/*
		fmt.Println("SENDING PROCEED FRAME BEFORE FILE")
		if err := sendFrame(c.PROCEED_FRAME, encoder); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
	*/

	fmt.Println("WAITING FOR PROCEED BEFORE FILE")
	if _, err := acceptFrame(decoder); err != nil {
		fmt.Printf("Error waiting for proceed file: %s\n", err.Error())
		return
	}

	if err := sendFileDataToClient(file, meta, connection, decoder); err != nil {
		fmt.Printf("Error while reading file data: %s\n", err.Error())
		if err := sendErrorFrame(encoder, err.Error()); err != nil {
			fmt.Printf("Error sending error frame: %s\n", err.Error())
		}
		return
	}

	fmt.Println("Successfully sent file data")
}

func openFileForTransfer(meta c.FileMetaData) (*os.File, error) {
	directoryName := meta.Owner
	filePath := fmt.Sprintf("%s/%s", directoryName, meta.Name)

	var file *os.File
	if _, err := os.Stat(directoryName); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(directoryName, 0777); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func sendFileDataToClient(file *os.File, meta c.FileMetaData, connection net.Conn, decoder *gob.Decoder) error {
	/*
		if meta.Size <= int64(c.MAX_BUFFER_SIZE) {
			buffer := make([]byte, meta.Size)
			if _, err := file.Read(buffer); err != nil {
				return err
			}
			if _, err := connection.Write(buffer); err != nil {
				return err
			}
			return nil
		}
	*/
	fileBuffer := make([]byte, c.TEMP_BUFFER_SIZE)
	totalBytesSent := 0
	for totalBytesSent <= int(meta.Size) {
		numOfBytes, err := file.Read(fileBuffer)
		if err != nil {
			if err == io.EOF {
				if numOfBytes > 0 {
					if _, err := connection.Write(fileBuffer[:numOfBytes]); err != nil {
						return err
					}
				}
				break
			} else {
				return err
			}
		}

		if _, err := connection.Write(fileBuffer[:numOfBytes]); err != nil {
			return err
		}
		totalBytesSent += numOfBytes
		fmt.Printf("Total Sent %d\n", totalBytesSent)
		fmt.Println("WAITING FOR PROCEED FRAME")
		if frame, err := acceptFrame(decoder); err != nil {
			return err
		} else if frame.Type != c.PROCEED_FRAME {
			return fmt.Errorf("Error: %s\n", err.Error())
		}
	}

	return nil
}

func acceptFileMetaData(connection net.Conn) (c.FileMetaData, error) {
	meta := c.FileMetaData{}
	decoder := gob.NewDecoder(connection)
	fmt.Println("WAITING FOR GOB")
	if err := decoder.Decode(&meta); err != nil {
		fmt.Printf("ERROR IN DECODER: %s\n", err.Error())
		return meta, err
	}

	return meta, nil
}
