package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	Id           string `bson:"_id" json:"id"`
	Username     string `bson:"username" json:"username"`
	Password     []byte `bson:"password" json:"password"`
	CreationDate string `bson:"creationDate" json:"createdAt"`
}

type Mongo struct {
	client  *mongo.Client
	users   *mongo.Collection
	records *mongo.Collection
	ctx     context.Context
}

func New() (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Printf("Error connecting to mongo: %s\n", err.Error())
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		fmt.Printf("Error pinging mongo %s\n", err.Error())
		return nil, err
	}

	users := client.Database("cloudStore").Collection("user")
	records := client.Database("cloudStore").Collection("records")
	return &Mongo{client: client, users: users, records: records}, nil
}

func (mongo *Mongo) GetUser(id string) (*User, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	singleResult := mongo.users.FindOne(mongo.ctx, filter)
	user := &User{}
	if err := singleResult.Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (mongo *Mongo) CreateUser(user *User) error {
	_, err := mongo.users.InsertOne(mongo.ctx, user)
	return err
}

func (mongo *Mongo) SearchUser(username, password string) ([]*User, error) {
	filter := bson.D{primitive.E{Key: "username", Value: username}}
	cursor, err := mongo.users.Find(mongo.ctx, filter)
	if err != nil {
		return nil, err
	}

	var results []*User
	if err = cursor.All(mongo.ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *Mongo) CreateRecord() (string, error) {
	fmt.Println("CREATED RECORD")
	return "", nil
}
