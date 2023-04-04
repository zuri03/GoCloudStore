package records

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/zuri03/GoCloudStore/common"
	m "go.mongodb.org/mongo-driver/mongo"
)

type DBCLient struct {
	client  *mongo.Client
	records *mongo.Collection
	ctx     context.Context
}

func NewDBClient(uri string) (*DBCLient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Printf("Error connecting to mongo: %s\n", err.Error())
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		fmt.Printf("Error pinging mongo %s\n", err.Error())
		return nil, err
	}

	records := client.Database("cloudStore").Collection("records")

	return &DBCLient{client: client, records: records}, nil
}

//If you cannot find record return no error with an empty record
func (mongo *DBCLient) GetRecord(key string) (*common.Record, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: key}}
	singleResult := mongo.records.FindOne(mongo.ctx, filter)
	record := &common.Record{}
	if err := singleResult.Decode(record); err != nil {
		if err == m.ErrNoDocuments {
			return &common.Record{}, nil
		}
		return nil, err
	}
	return record, nil
}

func (mongo *DBCLient) CreateRecord(record common.Record) error {
	_, err := mongo.records.InsertOne(mongo.ctx, record)
	if err != nil {
		return err
	}
	return nil
}

func (mongo *DBCLient) DeleteRecord(key string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: key}}
	_, err := mongo.records.DeleteOne(mongo.ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (mongo *DBCLient) ReplaceRecord(record *common.Record) error {
	filter := bson.D{primitive.E{Key: "_id", Value: record.Key}}
	singleResult := mongo.records.FindOneAndReplace(mongo.ctx, filter, record)
	return singleResult.Err()
}
