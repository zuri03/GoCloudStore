package records

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/zuri03/GoCloudStore/common"
)

type AllowedUserHandler struct {
	dbClient           recordDataBase
	routineTracker     *sync.WaitGroup
	paramsMiddleware   func(user, owner, key string) error
	resourceMiddleware func(key string, db recordDataBase) (common.Record, error)
	ownerMiddleware    func(id string, record common.Record) error
}

func (handler *AllowedUserHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	handler.routineTracker.Add(1)
	defer handler.routineTracker.Done()

	owner := req.FormValue("owner")
	key := req.FormValue("key")
	user := req.FormValue("user")

	if err := handler.paramsMiddleware(user, owner, key); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	/*
		if !authenticate(users, writer, req) {
			return
		}
	*/

	record, err := handler.resourceMiddleware(key, handler.dbClient)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := handler.ownerMiddleware(owner, record); err != nil {
		http.Error(writer, err.Error(), http.StatusForbidden)
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
