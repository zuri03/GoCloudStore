package users

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type CreateHandler struct {
	Users map[string]User
}

func (handler *CreateHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Bad Request"))
		return
	}
	username := req.FormValue("username")
	password := req.FormValue("password")

	if password == "" || username == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Username or password missing from request"))
		return
	}

	id := fmt.Sprintf("%s:%s", username, password)

	hash, err := bcrypt.GenerateFromPassword([]byte(id), bcryptCost)
	if err != nil {
		fmt.Println("ERROR GENERATING HASH")
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Internal Server Error"))
		return
	}
	//may hash the password as well

	now := time.Now()
	creationTime := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())

	user := User{
		Id:           hash,
		Username:     username,
		Password:     password,
		CreationDate: creationTime,
	}

	handler.Users[string(hash)] = user
}
