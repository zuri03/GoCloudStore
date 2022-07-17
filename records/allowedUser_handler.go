package records

import (
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

	//Middleware pipeline
	if !checkParamsRecords(writer, req) {
		return
	}

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

	if !resourceExists(key, handler.Keeper, writer) {
		return
	}

	if !checkOwner(owner, key, handler.Keeper, writer) {
		return
	}

	switch req.Method {
	case http.MethodPut:
		handler.Add(owner, user, key, writer)
	case http.MethodDelete:
		handler.Remove(owner, user, key, writer)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("Method not allowed"))
	}
}

func (handler *AllowedUserHandler) Add(owner, user, id string, writer http.ResponseWriter) {

}

func (handler *AllowedUserHandler) Remove(owner, user, id string, writer http.ResponseWriter) {

}
