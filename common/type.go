package common

//These types will be used when I refactor the client/storage communication
//for now they are not used anywhere
type MessageType int8

type PotocolMessage struct {
	Type MessageType
	Data []byte
}
