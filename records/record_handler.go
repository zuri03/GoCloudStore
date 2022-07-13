package records

import (
	"net/http"

	"github.com/zuri03/GoCloudStore/records/db"
)

type RecordHandler struct {
	Keeper *RecordKeeper
	Users  *db.Mongo
}

func (handler *RecordHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	if !checkParamsRecords(writer, req) {
		return
	}

	//autheticate the user here
	owner := req.FormValue("owner")
	key := req.FormValue("key")
	if req.Method == http.MethodPost {
		handler.CreateRecord(owner, key, writer)
		return
	}

	if !resourceExists(key, handler.Keeper, writer) {
		return
	}

	switch req.Method {
	case http.MethodGet:
		if !canView(owner, key, handler.Keeper, writer) {
			return
		}
		handler.GetRecord(key, writer)
	case http.MethodDelete:
		if !checkOwner(owner, key, handler.Keeper, writer) {
			return
		}
		handler.DeleteRecord(key, writer)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("Method not allowed"))
	}

}

func (handler *RecordHandler) CreateRecord(id, key string, writer http.ResponseWriter) {

}

func (handler *RecordHandler) GetRecord(key string, writer http.ResponseWriter) {

}

func (handler *RecordHandler) DeleteRecord(key string, writer http.ResponseWriter) {

}
