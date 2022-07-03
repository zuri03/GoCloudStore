package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"

	c "github.com/zuri03/GoCloudStore/common"
)

func acceptFrame(decoder *gob.Decoder) (c.ProtocolFrame, error) {
	var frame c.ProtocolFrame
	if err := decoder.Decode(&frame); err != nil {
		return frame, err
	}
	return frame, nil
}

func sendFrame(frameType c.FrameType, data []byte, encoder *gob.Encoder) error {
	frame := c.ProtocolFrame{
		Type:          frameType,
		PayloadLength: int64(len(data)),
		Data:          data,
	}
	if err := encoder.Encode(frame); err != nil {
		return err
	}
	return nil
}

func sendErrorFrame(encoder *gob.Encoder, message string) error {
	buffer := []byte(message)
	frame := c.ProtocolFrame{
		Type:          c.ERROR_FRAME,
		PayloadLength: int64(len(buffer)),
		Data:          buffer,
	}
	if err := encoder.Encode(frame); err != nil {
		return err
	}
	return nil
}

func decodeMetaData(frame c.ProtocolFrame) (FileMetaData, error) {
	ioBuffer := new(bytes.Buffer)
	fmt.Printf("Meta frame data lenth => %d\n", len(frame.Data))
	fmt.Printf("Meta frame data => %s\n", string(frame.Data))
	ioBuffer.Write(frame.Data)
	gob.Register(new(FileMetaData)) //May be using this function incorrectly
	decoder := gob.NewDecoder(ioBuffer)
	var meta FileMetaData
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
			if err := os.Mkdir(directoryName, 0644); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return nil, err
	}

	return file, nil
}
