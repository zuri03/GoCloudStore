package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
)

type DBClient struct {
	client  *mongo.Client
	users   *mongo.Collection
	records *mongo.Collection
	ctx     context.Context
}

func New() (*DBClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	//This line is used to connect to a mongodb k8s service
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))

	//This connection string is for docker compose
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		fmt.Printf("Error connecting to mongo: %s\n", err.Error())
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		fmt.Printf("Error pinging mongo %s\n", err.Error())
		return nil, err
	}

	records := client.Database("cloudStore").Collection("records")
	users := client.Database("cloudStore").Collection("user")

	return &DBClient{client: client, users: users, records: records}, nil
}
