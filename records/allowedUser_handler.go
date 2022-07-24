package records

import (
	"encoding/json"
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
	record, err := handler.dbClient.GetRecord(key)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	record.AllowedUsers = append(record.AllowedUsers, user)

	if err := handler.dbClient.ReplaceRecord(record); err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(writer)
	encoder.Encode(record)
}

func (handler *AllowedUserHandler) Remove(key, user string, writer http.ResponseWriter) {
	record, err := handler.dbClient.GetRecord(key)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	findKeyIndex := func(rec *common.Record) int {
		for idx, allowedUser := range rec.AllowedUsers {
			if allowedUser == user {
				return idx
			}
		}
		return -1
	}
	userIndex := findKeyIndex(record)
	if userIndex == -1 {
		http.Error(writer, fmt.Sprintf("Cannot find %s in allowed users", user), http.StatusNotFound)
		return
	}
	record.AllowedUsers = append(record.AllowedUsers[:userIndex], record.AllowedUsers[userIndex+1:]...)

	encoder := json.NewEncoder(writer)
	encoder.Encode(record)
}
