package records

import (
	"fmt"
	"net/http"
)

type DeleteHandler struct {
	Keeper *RecordKeeper
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

	if key == "" || password == "" || username == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Key, Username or password missing from request"))
		return
	}

	if err := handler.Keeper.RemoveRecord(key, username, password); err != nil {
		if err.Error() == "Unathorized" {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(fmt.Sprintf("%s is not athorized to delete this record", key)))
		} else {
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte(fmt.Sprintf("Error: record %s not found", key)))
		}
		return
	}

	writer.WriteHeader(http.StatusOK)
}
