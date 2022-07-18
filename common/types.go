package common

type FrameType uint8

type ProtocolFrame struct {
	Type          FrameType
	PayloadLength int64
	Data          []byte
}

type FileMetaData struct {
	Owner string
	Name  string
	Size  int64
}

type Record struct {
	Key          string   `bson:"_id" json:"key"`
	Size         int64    `bson:"size" json:"size"`
	Name         string   `bson:"name" json:"name"`
	Location     string   `bson:"location" json:"location"`
	CreatedAt    string   `bson:"createdAt" json:"createdAt"`
	IsPublic     bool     `bson:"isPublic" json:"isPublic"`
	Owner        string   `bson:"owner" json:"owner"`
	AllowedUsers []string `bson:"allowedUsers" json:"allowedUsers"`
}

type User struct {
	Id           string `bson:"_id" json:"id"`
	Username     string `bson:"username" json:"username"`
	Password     []byte `bson:"password" json:"password"`
	CreationDate string `bson:"creationDate" json:"createdAt"`
}
