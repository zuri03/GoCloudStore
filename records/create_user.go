package records

import (
	"encoding/json"
	"net/http"
)

type CreateHandler struct {
	Users *Users
}

func (handler *CreateHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	password := req.FormValue("password")

	user, err := handler.Users.create(username, password)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(user)
	writer.Write(jsonBytes)
}
