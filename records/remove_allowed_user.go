package records

import (
	"fmt"
	"net/http"
)

type RemoveUserHandler struct {
	Keeper *RecordKeeper
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

	if err := handler.Keeper.RemoveAllowedUser(key, username, password, removedUser); err != nil {
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
