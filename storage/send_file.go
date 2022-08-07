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

	fmt.Println("Sending proceed")

	if err := sendProceed(encoder); err != nil {
		fmt.Printf("Error sending proceed: %s\n", err.Error())
		return
	}

	if err := sendFileDataToClient(file, meta, connection, decoder); err != nil {
		fmt.Printf("%s\n", err.Error())
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

	if meta.Size <= int64(c.MAX_BUFFER_SIZE) {
		fmt.Println("SENDING FILE IN ONE CHUNCK")
		buffer := make([]byte, meta.Size)
		if _, err := file.Read(buffer); err != nil {
			return err
		}
		if _, err := connection.Write(buffer); err != nil {
			return err
		}
		return nil
	}

	fileBuffer := make([]byte, c.TEMP_BUFFER_SIZE)
	totalBytesSent := 0
	for totalBytesSent <= int(meta.Size) {
		numOfBytes, err := file.Read(fileBuffer)
		if err != nil && err != io.EOF {
			return err
		}

		if numOfBytes == 0 || err == io.EOF {
			fmt.Printf("GOT ALL OF THE BYTES BREAKING \n")
			//Just in case there are some bytes left over
			if numOfBytes != 0 {
				if _, err := connection.Write(fileBuffer[:numOfBytes]); err != nil {
					fmt.Printf("Error writing to final file: %s\n", err.Error())
					return err
				}
			}

			break
		}

		fmt.Printf("Sending %d bytes: %s\n", numOfBytes, string(fileBuffer[:numOfBytes]))

		if _, err := connection.Write(fileBuffer[:numOfBytes]); err != nil {
			fmt.Printf("Error while writing: %s\n", err.Error())
			return err
		}

		totalBytesSent += numOfBytes
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
