package records

import (
	"fmt"
	"net/http"
)

type DeleteHandler struct {
	Keeper *RecordKeeper
	Users  *Users
}

func (handler *DeleteHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Bad Request"))
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	key := req.FormValue("key")

	owner, err := handler.Users.get(username, password)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := handler.Keeper.RemoveRecord(key, owner.Id); err != nil {
		if err.Error() == "Unathorized" {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(fmt.Sprintf("%s is not athorized to delete this record", username)))
		} else {
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte(fmt.Sprintf("Error: record %s not found", key)))
		}
		return
	}

	writer.WriteHeader(http.StatusOK)
}
