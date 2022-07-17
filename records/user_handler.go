package records

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/zuri03/GoCloudStore/common"
)

type UserHandler struct {
	dbClient       Mongo
	routineTracker *sync.WaitGroup
}

func (handler *UserHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	handler.routineTracker.Add(1)
	defer handler.routineTracker.Done()

	encoder := json.NewEncoder(writer)
	if req.Method == http.MethodPost {

		if !checkParamsUsername(writer, req) {
			return
		}

		username := req.FormValue("username")
		password := req.FormValue("password")

		id, err := handler.CreateUser(username, password)

		if err != nil {
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		user, err := handler.GetUser(id)
		if err != nil {
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if err := encoder.Encode(*user); err != nil {
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		return
	}

	if !checkParamsId(writer, req) {
		return
	}

	if req.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := req.FormValue("id")

	user, err := handler.GetUser(id)
	if err != nil {
		http.Error(writer, "Internal Server Error, unable to create user", http.StatusInternalServerError)
		return
	}

	if err := encoder.Encode(*user); err != nil {
		fmt.Printf("Error encoding json: %s\n", err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (handler *UserHandler) CreateUser(username, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	now := time.Now().Format("2006-01-02 03:04:05")
	id := uuid.New()

	user := common.User{
		Id:           id.String(),
		Username:     username,
		Password:     hash,
		CreationDate: now,
	}

	err = handler.dbClient.CreateUser(&user)
	return id.String(), nil
}

func (handler *UserHandler) GetUser(id string) (*common.User, error) {
	user, err := handler.dbClient.GetUser(id)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Get user result => %+v", *user)
	return user, nil
}
