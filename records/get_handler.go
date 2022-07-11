package records

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GetHandler struct {
	Keeper *RecordKeeper
	Users  *Users
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

	owner, err := handler.Users.get(username, password)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	record, err := handler.Keeper.Get(key, owner.Id)
	if err != nil {
		fmt.Printf("Middleware has missed error: %s\n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Internal Server Error"))
		return
	}

	jsonBytes, _ := json.Marshal(record)
	writer.Write(jsonBytes)
}
