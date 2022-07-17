package db

import (
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
)

type Record struct {
	Key          string   `bson:"_id"`
	Size         int64    `bson:"size"`
	Name         string   `bson:"name"`
	Location     string   `bson:"location"`
	CreatedAt    string   `bson:"createdAt"`
	IsPublic     bool     `bson:"isPublic"`
	Owner        string   `bson:"owner"`
	AllowedUsers []string `bson:"allowedUsers"`
}

//If you cannot find record return no error with an empty record
func (mongo *Mongo) GetRecord(key, id string) (*Record, error) {
	return nil, nil
}

func (mongo *Mongo) CreateRecord(record *Record) error {
	return nil
}

func (mongo *Mongo) RemoveRecord(key, id string) error {
	return nil
}

func (mongo *Mongo) AddAllowedUser(key, id, user string) error {
	return nil
}

func (mongo *Mongo) RemoveAllowedUser(key, id, user string) error {
	return nil
}
