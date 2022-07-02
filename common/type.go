package common

//These types will be used when I refactor the client/storage communication
//for now they are not used anywhere
type FrameType int8

type ProtocolFrame struct {
	Type          FrameType
	PayloadLength int64
	Data          []byte
}
