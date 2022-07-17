package storage

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"

	c "github.com/zuri03/GoCloudStore/common"
)

//TODO: Implement connection reading with buffered io
type FileMetaData struct {
	Size     int64
	FileName string
	Username string
}

func storeFileHandler(connection net.Conn, frame c.ProtocolFrame) {
	fmt.Println("In store file handler")
	encoder, _ := newEncoderDecorder(connection)
	meta, err := decodeMetaData(frame)
	if err != nil {
		if err != io.EOF {
			sendErrorFrame(encoder, err.Error())
		}
		fmt.Printf("Error decoding meta: %s\n", err.Error())
		return
	}
	fmt.Printf("Get meta => %s \n size => %d\n", meta.FileName, meta.Size)

	if err := sendProceed(encoder); err != nil {
		if err != io.EOF {
			sendErrorFrame(encoder, err.Error())
		}
		fmt.Printf("Error sending proceed: %s\n", err.Error())
		return
	}

	if err := storeFileDataFromClient(meta, connection); err != nil {
		if err != io.EOF {
			sendErrorFrame(encoder, err.Error())
		}
		fmt.Printf("Error storing file data: %s\n", err.Error())
		return
	}

	fmt.Println("Successfully store file")
}

func storeFileDataFromClient(meta FileMetaData, connection net.Conn) error {

	file, err := openFile(meta.Username, meta.FileName)
	if err != nil {
		return err
	}
	defer file.Close()

	/*
		In order to fix the EOF errors the client must wait before it begins to send file data
		this write to the connection serves as a signal from the server to the client letting the client
		know that the server is ready to begin receiving file data
	*/
	//connection.Write([]byte(c.PROCEED_PROTOCOL))
	if meta.Size <= int64(c.MAX_CACHE_BUFFER_SIZE) {
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

	fileDataCacheBuffer := new(bytes.Buffer)
	readBuffer := make([]byte, c.TEMP_BUFFER_SIZE)
	writeBuffertoFile := func(buffer *bytes.Buffer, file *os.File) {
		file.Write(buffer.Bytes())
		buffer.Reset()
	}

	for {
		numOfBytes, err := connection.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				if numOfBytes > 0 {
					fileDataCacheBuffer.Write(readBuffer[:numOfBytes])
					writeBuffertoFile(fileDataCacheBuffer, file)
				}
				break
			}
			fmt.Println("ERROR OCCURRED RETURING")
			return err
		}

		if numOfBytes == 0 {
			fmt.Println("finished reading breaking")
			break
		}

		if _, err = fileDataCacheBuffer.Write(readBuffer[:numOfBytes]); err != nil {
			return err
		}

		if fileDataCacheBuffer.Len() > c.MAX_CACHE_BUFFER_SIZE {
			writeBuffertoFile(fileDataCacheBuffer, file)
		}
	}
	return nil
}
