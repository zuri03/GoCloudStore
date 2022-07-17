package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/zuri03/GoCloudStore/common"
)

func (mongo *Mongo) GetUser(id string) (*common.User, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	singleResult := mongo.users.FindOne(mongo.ctx, filter)
	user := &common.User{}
	if err := singleResult.Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (mongo *Mongo) CreateUser(user *common.User) error {
	_, err := mongo.users.InsertOne(mongo.ctx, user)
	return err
}

func (mongo *Mongo) SearchUser(username, password string) ([]*common.User, error) {
	filter := bson.D{primitive.E{Key: "username", Value: username}}
	cursor, err := mongo.users.Find(mongo.ctx, filter)
	if err != nil {
		return nil, err
	}

	var results []*common.User
	if err = cursor.All(mongo.ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
