package records

import (
	"net/http"
)

type UserClient http.Client

type User struct {
	Id [32]byte
}

func (u *UserClient) Get(username string, password string) (User, error) {
	return User{}, nil
}

func (u *UserClient) Exists(username string, password string) (bool, error) {
	return true, nil
}
