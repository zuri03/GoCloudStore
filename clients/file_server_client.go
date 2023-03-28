package client

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/zuri03/GoCloudStore/common"
)

type FileServerclient struct {
	Client http.Client
}

func NewFileServerClient() *FileServerclient {
	return nil
}

func newEncoderDecoder(connection net.Conn) (*gob.Encoder, *gob.Decoder) {
	gob.Register(new(common.ProtocolFrame))
	gob.Register(new(common.FileMetaData))

	return gob.NewEncoder(connection), gob.NewDecoder(connection)
}

func (client *FileServerclient) SendFile(owner string, file *os.File, fileInfo fs.FileInfo) error {
	//TODO: The address of the datanode must come from the record server
	connection, err := net.DialTimeout("tcp", ":8000", time.Duration(10)*time.Second)
	if err != nil {
		return err
	}
	defer connection.Close()

	encoder, decoder := newEncoderDecoder(connection)
	/*
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
	*/

	meta := common.FileMetaData{
		Owner: owner,
		Name:  fileInfo.Name(),
		Size:  fileInfo.Size(),
	}

	if err := sendMetaDataToServer(common.SEND_FRAME, meta, encoder); err != nil {
		fmt.Printf("Error sending meta data to server: %s\n", err.Error())
		return err
	}

	if _, err := acceptFrame(decoder, common.PROCEED_FRAME); err != nil {
		fmt.Printf("Error waiting for proceed: %s\n", err.Error())
		return err
	}

	if err := sendFileDataToServer(file, meta, connection, decoder); err != nil {
		fmt.Printf("Error sending file data to server: %s\n", err.Error())
		return err
	}

	fmt.Println("Successfully sent file to server")
	return nil
}

/*
*Retrieves a file from a server with ID fileId
	- fileId: unique identifier of file that needs to be retrieved from server

Returns:
	- string: file path pointing to tmp file downloaded from the server and saved on this machine
*/
//TODO: Return the file path to the temp file that was just downloaded
func (client *FileServerclient) GetFile(owner string, record *common.Record) (string, error) {
	connection, err := net.DialTimeout("tcp", ":8000", time.Duration(10)*time.Second)
	if err != nil {
		return "", err
	}
	defer connection.Close()
	meta := common.FileMetaData{
		Owner: owner,
		Name:  record.Name,
		Size:  record.Size,
	}
	encoder, decoder := newEncoderDecoder(connection)
	if err := sendMetaDataToServer(common.GET_FRAME, meta, encoder); err != nil {
		return "", err
	}

	fmt.Println("Waiting for proceed")

	if _, err := acceptFrame(decoder, common.PROCEED_FRAME); err != nil {
		return "", err
	}

	if err := getFileDataFromServer(meta.Name, int(meta.Size), connection, encoder); err != nil {
		return "", err
	}

	return "", nil
}

//uses a connection to retrieve byte data from a storage server and store it in a file
func getFileDataFromServer(filePath string, fileSize int, connection net.Conn, encoder *gob.Encoder) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return err
	}
	defer file.Close()

	if fileSize <= common.MAX_BUFFER_SIZE {
		fmt.Println("Storing file in one chunk")
		buffer := make([]byte, fileSize)
		numOfBytes, err := connection.Read(buffer)
		fmt.Printf("Buffer: %s\n", buffer)
		if err != nil {
			return err
		}
		if _, err := file.Write(buffer[:numOfBytes]); err != nil {
			return err
		}
		return nil
	}

	readBuffer := make([]byte, common.TEMP_BUFFER_SIZE)
	bufferedWriter := bufio.NewWriterSize(file, common.MAX_BUFFER_SIZE)
	defer bufferedWriter.Flush()
	fmt.Println("SENDING IN LOOP")
	totalBytes := 0
	for totalBytes <= fileSize {
		fmt.Println("READY TO READ")
		numOfBytes, err := connection.Read(readBuffer)
		if err != nil && err != io.EOF {
			return err
		}

		if numOfBytes == 0 || err == io.EOF {
			fmt.Printf("GOT ALL OF THE BYTES BREAKING \n")
			//Just in case there are some bytes left over
			if numOfBytes != 0 {
				if _, err := bufferedWriter.Write(readBuffer[:numOfBytes]); err != nil {
					fmt.Printf("Error writing to final file: %s\n", err.Error())
					return err
				}
			}

			break
		}

		fmt.Printf("read %d bytes from connection: %s\n", numOfBytes, string(readBuffer[:numOfBytes]))
		if _, err := bufferedWriter.Write(readBuffer[:numOfBytes]); err != nil {
			fmt.Printf("error writing to file: %s\n", err.Error())
			return err
		}

		totalBytes += numOfBytes
	}

	return nil
}

func sendFileDataToServer(file *os.File, meta common.FileMetaData, connection net.Conn, decoder *gob.Decoder) error {
	/*
		if meta.Size <= int64(common.MAX_BUFFER_SIZE) {
			fmt.Println("Sending file in one chunck")
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
	fmt.Println("SENDING FILE DATA")
	buffer := make([]byte, common.TEMP_BUFFER_SIZE)
	for {
		numOfBytes, err := file.Read(buffer)

		if err != nil {
			if err == io.EOF {
				if numOfBytes > 0 {
					if _, err := connection.Write(buffer[:numOfBytes]); err != nil {
						return err
					}
				}
				return nil
			} else {
				return err
			}
		}

		if _, err := connection.Write(buffer[:numOfBytes]); err != nil {
			return err
		}
		fmt.Println("WAITING FOR PROCEED")
		if _, err := acceptFrame(decoder, common.PROCEED_FRAME); err != nil {
			return err
		}
	}
}

func sendFrame(frameType common.FrameType, encoder *gob.Encoder, data ...[]byte) error {

	var frame common.ProtocolFrame
	if len(data) > 0 {
		frame = common.ProtocolFrame{
			Type:          frameType,
			PayloadLength: int64(len(data[0])),
			Data:          data[0],
		}
	} else {
		frame = common.ProtocolFrame{
			Type:          frameType,
			PayloadLength: 0,
			Data:          nil,
		}
	}

	if err := encoder.Encode(frame); err != nil {
		return err
	}

	return nil
}

func acceptFrame(decoder *gob.Decoder, acceptedTypes ...common.FrameType) (*common.ProtocolFrame, error) {
	var frame common.ProtocolFrame
	if err := decoder.Decode(&frame); err != nil {
		return nil, err
	}

	if frame.Type == common.ERROR_FRAME {
		return nil, fmt.Errorf("Server Error: %s\n", string(frame.Data))
	}

	for _, frameType := range acceptedTypes {
		if frameType == frame.Type {
			return &frame, nil
		}
	}

	return nil, fmt.Errorf("Unexpected frame: %d\n", frame.Type)
}

func sendMetaDataToServer(frameType common.FrameType, meta common.FileMetaData, encoder *gob.Encoder) error {
	metaBuffer := new(bytes.Buffer)
	if err := gob.NewEncoder(metaBuffer).Encode(meta); err != nil {
		return err
	}

	frame := common.ProtocolFrame{
		Type:          frameType,
		PayloadLength: int64(metaBuffer.Len()),
		Data:          metaBuffer.Bytes(),
	}

	return encoder.Encode(frame)
}
