package db

import (
	"github.com/zuri03/GoCloudStore/common"
)

//If you cannot find record return no error with an empty record
func (mongo *Mongo) GetRecord(key string) (*common.Record, error) {
	return nil, nil
}

func (mongo *Mongo) CreateRecord(record common.Record) error {
	return nil
}

func (mongo *Mongo) DeleteRecord(key string) error {
	return nil
}

func (mongo *Mongo) AddAllowedUser(key, user string) error {
	return nil
}

func (mongo *Mongo) RemoveAllowedUser(key, user string) error {
	return nil
}
