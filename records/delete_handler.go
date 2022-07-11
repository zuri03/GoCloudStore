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
		fmt.Printf("Middleware did not catch error: %s\n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := handler.Keeper.Remove(key, owner.Id); err != nil {
		fmt.Printf("Middleware did not catch error: %s\n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
