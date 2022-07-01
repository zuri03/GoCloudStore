package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	sha "crypto/sha256"
)

type GetHandler struct {
	Users map[[32]byte]User
}

func (handler *GetHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("Got GET request")
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

	id := []byte(fmt.Sprintf("%s:%s", username, password))
	hash := sha.Sum256(id)
	user, ok := handler.Users[hash]
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
