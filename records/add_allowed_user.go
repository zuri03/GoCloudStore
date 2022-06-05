package records

import (
	"fmt"
	"net/http"
)

type AddUserHandler struct {
	Keeper *RecordKeeper
}

func (handler *AddUserHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Bad Request"))
		return
	}

	allowedUser := req.FormValue("allowedUser")
	username := req.FormValue("username")
	password := req.FormValue("password")
	key := req.FormValue("key")

	if key == "" || password == "" || username == "" || allowedUser == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Key, Username, password or allowed user missing from request"))
		return
	}

	if err := handler.Keeper.AddAllowedUser(key, username, password, allowedUser); err != nil {
		if err.Error() == "Unathorized" {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(fmt.Sprintf("%s is not athorized to add allowed users to this record", key)))
		} else {
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte(fmt.Sprintf("Error: record %s not found", key)))
		}
		return
	}

	writer.WriteHeader(http.StatusOK)
}
