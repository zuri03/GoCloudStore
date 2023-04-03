package user

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UserService struct {
	db *DBClient
}

func (userService *UserService) CreateUser(username string, password string) ([]byte, int, error) {
	return nil, http.StatusOK, nil
}

func (userService *UserService) GetUser(id string) ([]byte, int, error) {
	if id == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("Error: ID is missing from request")
	}

	user, err := userService.db.GetUser(id)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if user == nil {
		return nil, http.StatusNotFound, fmt.Errorf("Unable to find user with ID %s\n", id)
	}

	jsonBytes, err := json.Marshal(user)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return jsonBytes, http.StatusOK, nil
}

func (userService *UserService) DeleteUser(id string) (int, error) {
	return http.StatusOK, nil
}

func (userService *UserService) AuthorizeUser(username string, password string) ([]byte, int, error) {
	return nil, http.StatusOK, nil
}
