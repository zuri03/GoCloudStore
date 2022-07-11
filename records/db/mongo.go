package db

import (
	_"context"
	_"go.mongodb.org/mongo-driver/mongo"
    _"go.mongodb.org/mongo-driver/mongo/options"
    _"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongo struct {
	Client Mongo.Client
}

func (m *Mongo) New() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        return err
    }
}