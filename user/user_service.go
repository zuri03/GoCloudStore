package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zuri03/GoCloudStore/common"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *DBClient
}

func (userService *UserService) CreateUser(username string, password string) ([]byte, int, error) {
	if username == "" || password == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("username or password missing from request got, username: %s, password: %s\n", username, password)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	now := time.Now().Format("2006-01-02 03:04:05")
	id := uuid.New()

	newUser := common.User{
		Id:           id.String(),
		Username:     username,
		Password:     hash,
		CreationDate: now,
	}

	if err := userService.db.CreateUser(&newUser); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	jsonBytes, err := json.Marshal(newUser)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return jsonBytes, http.StatusOK, nil
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
	if id == "" {
		return http.StatusBadRequest, fmt.Errorf("Error: ID is missing from request")
	}

	deletedCount, err := userService.db.DeleteUser(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if deletedCount == 0 {
		return http.StatusNotFound, fmt.Errorf("Error: unable to find user with ID %s\n", id)
	}

	return http.StatusOK, nil
}

type Response struct {
	Id string `json:"id"`
}

func (userService *UserService) AuthorizeUser(username string, password string) ([]byte, int, error) {
	if username == "" || password == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("username or password missing from request got, username: %s, password: %s\n", username, password)
	}

	potentialUsers, err := userService.db.SearchUser(username, password)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	for _, potentialMatch := range potentialUsers {
		matchPassword := potentialMatch.Password
		err := bcrypt.CompareHashAndPassword(matchPassword, []byte(password))
		if err == nil {
			response := Response{Id: potentialMatch.Id}
			jsonBytes, err := json.Marshal(response)
			if err != nil {
				return nil, http.StatusInternalServerError, err
			}

			return jsonBytes, http.StatusOK, nil
		}
	}

	return nil, http.StatusNotFound, fmt.Errorf("Error: unable to find user with username %s and password %s\n", username, password)
}
