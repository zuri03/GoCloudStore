package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"

	c "github.com/zuri03/GoCloudStore/common"
)

func newEncoderDecorder(connection net.Conn) (*gob.Encoder, *gob.Decoder) {
	gob.Register(new(c.ProtocolFrame))
	gob.Register(new(c.FileMetaData))

	return gob.NewEncoder(connection), gob.NewDecoder(connection)
}

func acceptFrame(decoder *gob.Decoder) (c.ProtocolFrame, error) {
	var frame c.ProtocolFrame
	if err := decoder.Decode(&frame); err != nil {
		return frame, err
	}
	return frame, nil
}

//A general method used from sending frames
func sendFrame(frameType c.FrameType, encoder *gob.Encoder, data ...[]byte) error {
	if len(data) == 0 {
		frame := c.ProtocolFrame{
			Type:          frameType,
			PayloadLength: 0,
			Data:          nil,
		}

		return encoder.Encode(frame)
	}

	for _, payload := range data {
		frame := c.ProtocolFrame{
			Type:          frameType,
			PayloadLength: int64(len(payload)),
			Data:          payload,
		}

		if err := encoder.Encode(frame); err != nil {
			return err
		}
	}

	return nil
}

//A helper method to simplify the numerous attempts to send an error frame to the client
func sendErrorFrame(encoder *gob.Encoder, message string) error {
	buffer := []byte(message)
	frame := c.ProtocolFrame{
		Type:          c.ERROR_FRAME,
		PayloadLength: int64(len(buffer)),
		Data:          buffer,
	}
	if err := encoder.Encode(frame); err != nil {
		if err != io.EOF {
			return err
		}
		fmt.Printf("EOF ERROR ON ERROR FRAME: %s\n", err.Error())
	}
	return nil
}

func sendSuccessFrame(encoder *gob.Encoder) error {
	frame := c.ProtocolFrame{
		Type:          c.SUCCESS_FRAME,
		PayloadLength: 0,
		Data:          nil,
	}

	if err := encoder.Encode(frame); err != nil {
		if err != io.EOF {
			return err
		}
		fmt.Printf("EOF ERROR ON SUCCESS FRAME: %s\n", err.Error())
	}

	return nil
}

func decodeMetaData(frame c.ProtocolFrame) (c.FileMetaData, error) {
	fmt.Println("DECODING META DATA")
	ioBuffer := new(bytes.Buffer)
	ioBuffer.Write(frame.Data)
	var meta c.FileMetaData
	decoder := gob.NewDecoder(ioBuffer)
	if err := decoder.Decode(&meta); err != nil {
		return meta, err
	}
	return meta, nil
}

func openFile(directoryName string, fileName string) (*os.File, error) {
	filePath := fmt.Sprintf("%s/%s", directoryName, fileName)

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

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func sendProceed(encoder *gob.Encoder) error {
	proceedFrame := c.ProtocolFrame{
		Type:          c.PROCEED_FRAME,
		PayloadLength: 0,
		Data:          nil,
	}

	if err := encoder.Encode(proceedFrame); err != nil {
		return err
	}

	return nil
}
