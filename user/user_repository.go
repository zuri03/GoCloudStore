package user

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
)

type DBClient struct {
	client  *mongo.Client
	users   *mongo.Collection
	records *mongo.Collection
	ctx     context.Context
}

func NewDBClient(uri string) (*DBClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	//This line is used to connect to a mongodb k8s service
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))

	//This connection string is for docker compose
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
	users := client.Database("cloudStore").Collection("user")

	return &DBClient{client: client, users: users, records: records}, nil
}

func (dbClient *DBClient) GetUser(id string) (*common.User, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	singleResult := dbClient.users.FindOne(dbClient.ctx, filter)
	user := &common.User{}
	if err := singleResult.Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return user, nil
}

func (dbClient *DBClient) CreateUser(user *common.User) error {
	_, err := dbClient.users.InsertOne(dbClient.ctx, user)
	return err
}

func (dbClient *DBClient) SearchUser(username, password string) ([]*common.User, error) {
	filter := bson.D{primitive.E{Key: "username", Value: username}}
	cursor, err := dbClient.users.Find(dbClient.ctx, filter)
	if err != nil {
		return nil, err
	}

	var results []*common.User
	if err = cursor.All(dbClient.ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (dbClient *DBClient) DeleteUser(id string) (int64, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	deleteResult, err := dbClient.users.DeleteOne(dbClient.ctx, filter)
	if err != nil {
		return 0, err
	}
	return deleteResult.DeletedCount, nil
}
