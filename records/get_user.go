package records

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GetUserHandler struct {
	Users *Users
}

func (handler *GetUserHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	password := req.FormValue("password")

	user, err := handler.Users.get(username, password)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user == nil {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(fmt.Sprintf("User %s not found", username)))
		return
	}

	jsonBytes, err := json.Marshal(user)
	writer.Write(jsonBytes)
}
