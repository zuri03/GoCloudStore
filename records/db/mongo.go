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
	client *mongo.Client
	users  *mongo.Collection
	ctx    context.Context
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

	collection := client.Database("cloudStore").Collection("user")
	return &Mongo{client: client, users: collection}, nil
}

func (m *Mongo) GetUser(id string) (*User, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	singleResult := m.users.FindOne(m.ctx, filter)
	user := &User{}
	if err := singleResult.Decode(user); err != nil {
		return nil, err
	}
	return user, nil
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

	fmt.Printf("Results => %+v\n", results)
	fmt.Println("pretty print")

	for _, c := range results {
		fmt.Printf("result => %+v\n", c)
	}

	return results, nil
}
