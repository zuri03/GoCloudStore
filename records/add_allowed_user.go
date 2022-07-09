package records

import (
	"fmt"
	"net/http"
	"strings"
)

type AddUserHandler struct {
	Keeper *RecordKeeper
	Users  *Users
}

func (handler *AddUserHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("IN ADD HANDLER")
	if err := req.ParseForm(); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Bad Request"))
		return
	}

	allowedUser := req.FormValue("allowedUser")
	username := req.FormValue("username")
	password := req.FormValue("password")
	key := req.FormValue("key")

	if allowedUser == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Key, Username, password or allowed user missing from request"))
		return
	}

	allowedUserCreds := strings.Split(allowedUser, ":")

	owner, err := handler.Users.get(username, password)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	allowed, err := handler.Users.get(allowedUserCreds[0], allowedUserCreds[1])
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if allowed == nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf("Cannot find user: %s", allowedUserCreds[0])))
		return
	}

	if err := handler.Keeper.AddAllowedUser(key, owner.Id, allowed.Id); err != nil {
		if err.Error() == "Unathorized" {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(fmt.Sprintf("%s is not athorized to add allowed users to this record", username)))
		} else {
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte(fmt.Sprintf("Error: record %s not found", key)))
		}
		return
	}

	writer.WriteHeader(http.StatusOK)
}
