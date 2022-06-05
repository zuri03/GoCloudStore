package records

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GetHandler struct {
	Keeper *RecordKeeper
}

func (handler *GetHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
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

	record, err := handler.Keeper.GetRecord(key, username, password)
	if err != nil {
		if err.Error() == "Unathorized" {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(fmt.Sprintf("%s is not athorized to view this record", key)))
		} else {
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte(fmt.Sprintf("Error: record %s not found", key)))
		}
		return
	}

	jsonBytes, _ := json.Marshal(record)
	writer.Write(jsonBytes)
}
