package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	Id           string `bson:"_id"`
	Username     string `bson:"username"`
	Password     []byte `bson:"password"`
	CreationDate string `bson:"creationDate"`
}

type Mongo struct {
	client *mongo.Client
	users  *mongo.Collection
	ctx    context.Context
}

func New() (*Mongo, error) {
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	collection := client.Database("cloudStore").Collection("user")
	return &Mongo{client: client, users: collection}, nil
}

func (m *Mongo) CreateUser(user *User) error {
	_, err := m.users.InsertOne(m.ctx, user)
	return err
}

func (m *Mongo) SearchUser(username, password string) ([]*User, error) {
	filter := bson.D{primitive.E{Key: "username", Value: username}}
	cursor, err := m.users.Find(m.ctx, filter)
	if err != nil {
		return nil, err
	}

	var results []*User
	if err = cursor.All(m.ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
