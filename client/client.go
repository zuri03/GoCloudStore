package client

import (
	"github.com/zuri03/GoCloudStore/common"
)

type client struct {
}

func New() *client {
	return nil
}

func (client *client) SendFile(filePath string) error {
	return nil
}

/*
*Retrieves a file from a server with ID fileId
	- fileId: unique identifier of file that needs to be retrieved from server

Returns:
	- string: file path pointing to tmp file downloaded from the server and saved on this machine
*/
func (client *client) GetFile(fileId string) (string, error) {
	return "", nil
}

func (client *client) AuthenticateUser(username string, password string) (string, error) {
	return "", nil
}

//returns ID of newly created user
func (client *client) CreateUser(username string, password string) (string, error) {
	return "", nil
}

func (client *client) GetFileRecord(owner string, key string) (*common.Record, error) {
	return nil, nil
}

//returns error in case file record was unable to be deleted
func (client *client) DeleteFileRecord(owner string, key string) error {
	return nil
}

func (client *client) CreateFileRecord(owner, key, fileName string, fileSize int64) error {
	return nil
}

func (client *client) AddAllowedUser(owner, key, allowedUser string) error {
	return nil
}

func (client *client) RemoveAllowedUser(owner, key string, removedUser string) error {
	return nil
}

func (client *client) sendHttpRequest(url string) {

}
