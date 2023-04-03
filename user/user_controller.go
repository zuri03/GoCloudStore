package user

import (
	"net/http"
)

func CreateUser(username string, password string) ([]byte, int, error) {
	return nil, http.StatusOK, nil
}

func GetUser(id string) ([]byte, int, error) {
	return nil, http.StatusOK, nil
}

func DeleteUser(id string) (int, error) {
	return http.StatusOK, nil
}

func AuthorizeUser(username string, password string) ([]byte, int, error) {
	return nil, http.StatusOK, nil
}
