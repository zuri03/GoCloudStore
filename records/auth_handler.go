package records

import (
	"fmt"
	"net/http"
)

type AuthHandler struct {
	Users *Users
}

func (handler *AuthHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	password := req.FormValue("password")

	exist, err := handler.Users.exists(username, password)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if exist {
		fmt.Println("User exist")
		writer.WriteHeader(http.StatusOK)
		return
	} else {
		fmt.Println("User does not exist")
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
}
