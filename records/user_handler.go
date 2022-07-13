package records

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zuri03/GoCloudStore/records/db"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	users *db.Mongo
}

func (handler *UserHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		if !checkParamsUsername(writer, req) {
			return
		}

		username := req.FormValue("username")
		password := req.FormValue("password")

		if err := handler.CreateUser(username, password); err != nil {
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
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

	jsonBytes, err := json.Marshal(*user)
	if err != nil {
		http.Error(writer, "Internal Server Error, unable to marshal user", http.StatusInternalServerError)
		return
	}

	writer.Write(jsonBytes)
}

func (handler *UserHandler) CreateUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now().Format("2006-01-02 03:04:05")
	id := uuid.New()

	user := db.User{
		Id:           id.String(),
		Username:     username,
		Password:     hash,
		CreationDate: now,
	}

	err = handler.users.CreateUser(&user)
	return nil
}

func (handler *UserHandler) GetUser(id string) (*db.User, error) {
	return &db.User{}, nil
}
