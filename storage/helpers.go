package storage

import (
	"encoding/gob"

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
