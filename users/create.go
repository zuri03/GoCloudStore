package users

import (
	"fmt"
	"net/http"
	"time"

	sha "crypto/sha256"
)

type CreateHandler struct {
	Users map[[32]byte]User
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

	id := []byte(fmt.Sprintf("%s:%s", username, password))
	hash := sha.Sum256(id)
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

	handler.Users[hash] = user
}
