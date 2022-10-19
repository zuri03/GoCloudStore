package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/zuri03/GoCloudStore/common"
)

func (dbClient *DBClient) GetUser(id string) (*common.User, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	singleResult := dbClient.users.FindOne(dbClient.ctx, filter)
	user := &common.User{}
	if err := singleResult.Decode(user); err != nil {
		return nil, err
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
