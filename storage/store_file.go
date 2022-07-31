package storage

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/zuri03/GoCloudStore/common"
	c "github.com/zuri03/GoCloudStore/common"
)

//TODO: Implement connection reading with buffered io
func storeFileHandler(connection net.Conn, frame c.ProtocolFrame) {
	fmt.Println("In store file handler")
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

	if err := sendProceed(encoder); err != nil {
		if err != io.EOF {
			if err := sendErrorFrame(encoder, err.Error()); err != nil {
				fmt.Printf("Error sending error frame: %s\n", err.Error())
			}
		}
		fmt.Printf("Error sending proceed: %s\n", err.Error())
		return
	}

	if err := storeFileDataFromClient(meta, connection); err != nil {
		if err != io.EOF {
			if err := sendErrorFrame(encoder, err.Error()); err != nil {
				fmt.Printf("Error sending error frame: %s\n", err.Error())
			}
		}
		fmt.Printf("Error storing file data: %s\n", err.Error())
		return
	}

	fmt.Println("Successfully store file")
}

func storeFileDataFromClient(meta c.FileMetaData, connection net.Conn) error {

	file, err := openFile(meta.Owner, meta.Name)
	if err != nil {
		return err
	}
	defer file.Close()

	if meta.Size <= int64(c.MAX_BUFFER_SIZE) {
		fmt.Println("storing file in one chunk")
		buffer := make([]byte, meta.Size)
		if _, err := connection.Read(buffer); err != nil {
			return err
		}
		if _, err := file.Write(buffer); err != nil {
			return err
		}
		return nil
	}

	readBuffer := make([]byte, c.TEMP_BUFFER_SIZE)
	bufferedWriter := bufio.NewWriterSize(file, common.MAX_BUFFER_SIZE)

	for {
		numOfBytes, err := connection.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				if numOfBytes > 0 {
					if _, err := bufferedWriter.Write(readBuffer[:numOfBytes]); err != nil {
						return err
					}
				}
				break
			}
			fmt.Println("ERROR OCCURRED RETURING")
			return err
		}

		if _, err = bufferedWriter.Write(readBuffer[:numOfBytes]); err != nil {
			return err
		}
	}
	return nil
}
