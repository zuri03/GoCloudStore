package records

import (
	"net/http"

	"github.com/zuri03/GoCloudStore/records/db"
)

type AllowedUserHandler struct {
	Keeper *RecordKeeper
	Users  *db.Mongo
}

func (handler *AllowedUserHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
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
