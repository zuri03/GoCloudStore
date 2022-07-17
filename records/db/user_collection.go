package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	Id           string `bson:"_id" json:"id"`
	Username     string `bson:"username" json:"username"`
	Password     []byte `bson:"password" json:"password"`
	CreationDate string `bson:"creationDate" json:"createdAt"`
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
