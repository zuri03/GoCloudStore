package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/zuri03/GoCloudStore/common"
)

//If you cannot find record return no error with an empty record
func (mongo *Mongo) GetRecord(key string) (*common.Record, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: key}}
	singleResult := mongo.records.FindOne(mongo.ctx, filter)
	record := &common.Record{}
	if err := singleResult.Decode(record); err != nil {
		return nil, err
	}
	return nil, nil
}

func (mongo *Mongo) CreateRecord(record common.Record) error {
	_, err := mongo.records.InsertOne(mongo.ctx, record)
	if err != nil {
		return err
	}
	return nil
}

func (mongo *Mongo) DeleteRecord(key string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: key}}
	_, err := mongo.records.DeleteOne(mongo.ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (mongo *Mongo) ReplaceRecord(record *common.Record) error {
	filter := bson.D{primitive.E{Key: "_id", Value: record.Key}}
	singleResult := mongo.records.FindOneAndReplace(mongo.ctx, filter, record)
	return singleResult.Err()
}
