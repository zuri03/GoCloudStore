package records

import (
	"fmt"
	"net/http"
	"strings"
)

type RemoveUserHandler struct {
	Keeper *RecordKeeper
	Users  *Users
}

func (handler *RemoveUserHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("IN REMOVE HANDLER")
	if err := req.ParseForm(); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Bad Request"))
		return
	}

	removedUser := req.FormValue("removedUser")
	username := req.FormValue("username")
	password := req.FormValue("password")
	key := req.FormValue("key")

	if removedUser == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Key, Username, password or allowed user missing from request"))
		return
	}

	removedUserCreds := strings.Split(removedUser, ":")

	owner, err := handler.Users.get(username, password)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	removed, err := handler.Users.get(removedUserCreds[0], removedUserCreds[1])
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if removed == nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf("Cannot find user: %s", removedUserCreds[0])))
		return
	}

	if err := handler.Keeper.RemoveAllowedUser(key, owner.Id, removed.Id); err != nil {
		if err.Error() == "Unathorized" {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(fmt.Sprintf("%s is not athorized to add allowed users to this record", username)))
		} else {
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte(err.Error()))
		}
		return
	}

	writer.WriteHeader(http.StatusOK)
}
