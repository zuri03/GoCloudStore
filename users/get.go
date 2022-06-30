package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type GetHandler struct {
	Users map[string]User
}

func (handler *GetHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	fmt.Printf("Got GET request")
	if err := req.ParseForm(); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Bad Request"))
		return
	}
	username := req.FormValue("username")
	password := req.FormValue("password")
	fmt.Printf("Username => %s\n Password => %s\n", username, password)
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

	user, ok := handler.Users[string(hash)]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Write(jsonBytes)
}
