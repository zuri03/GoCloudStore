package records

import (
	"fmt"
	"net/http"
	"sync"
)

type AllowedUserHandler struct {
	dbClient       Mongo
	routineTracker *sync.WaitGroup
}

func (handler *AllowedUserHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	handler.routineTracker.Add(1)
	defer handler.routineTracker.Done()

	if !checkParamsAllowedUser(writer, req) {
		return
	}
	/*
		if !authenticate(users, writer, req) {
			return
		}
	*/
	owner := req.FormValue("owner")
	key := req.FormValue("key")
	user := req.FormValue("user")

	record, ok := resourceExists(key, handler.dbClient, writer)
	if !ok {
		return
	}

	if !checkOwner(owner, record, writer) {
		return
	}

	switch req.Method {
	case http.MethodPut:
		handler.Add(key, user, writer)
	case http.MethodDelete:
		handler.Remove(key, user, writer)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("Method not allowed"))
	}
}

func (handler *AllowedUserHandler) Add(key, user string, writer http.ResponseWriter) {
	if err := handler.dbClient.AddAllowedUser(key, user); err != nil {
		fmt.Printf("Error adding allowed user: %s\n", err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (handler *AllowedUserHandler) Remove(key, user string, writer http.ResponseWriter) {
	if err := handler.dbClient.RemoveAllowedUser(key, user); err != nil {
		fmt.Printf("Error adding allowed user: %s\n", err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
