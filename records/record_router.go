package records

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type RecordCreationRequest struct {
	Owner    string `json:"owner"`
	Key      string `json:"key"`
	FileName string `json:"name"`
	Size     int    `json:"size"`
}

type RecordRouter struct {
	waitgroup     *sync.WaitGroup
	recordService *RecordService
}

func (recordRouter *RecordRouter) ServeRecordHTTP(writer http.ResponseWriter, req *http.Request) {
	recordRouter.waitgroup.Add(1)
	defer recordRouter.waitgroup.Done()

	switch req.Method {
	case http.MethodPost:
		var requestBody RecordCreationRequest
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = json.Unmarshal(body, &requestBody); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		response, statusCode, err := recordRouter.recordService.CreateRecord(requestBody)
		if err != nil {
			http.Error(writer, err.Error(), statusCode)
			return
		}

		writer.WriteHeader(statusCode)
		writer.Write(response)
	case http.MethodGet:
		id := req.FormValue("id")
		key := req.FormValue("key")

		response, statusCode, err := recordRouter.recordService.GetRecord(id, key)
		if err != nil {
			http.Error(writer, err.Error(), statusCode)
			return
		}

		writer.WriteHeader(statusCode)
		writer.Write(response)
	case http.MethodDelete:
		id := req.FormValue("id")
		key := req.FormValue("key")

		statusCode, err := recordRouter.recordService.DeleteRecord(id, key)
		if err != nil {
			http.Error(writer, err.Error(), statusCode)
			return
		}

		writer.WriteHeader(statusCode)
	default:
		http.Error(writer, fmt.Sprintf("Method %s is not supported on this endpoint", req.Method), http.StatusMethodNotAllowed)
	}
}

func (recordRouter *RecordRouter) ServePrivilegedUserHTTP(writer http.ResponseWriter, req *http.Request) {
	recordRouter.waitgroup.Add(1)
	defer recordRouter.waitgroup.Done()

	owner := req.FormValue("owner")
	key := req.FormValue("key")
	user := req.FormValue("user")

	switch req.Method {
	case http.MethodPut:
		response, statusCode, err := recordRouter.recordService.AddPrivilegedUser(owner, key, user)
		if err != nil {
			http.Error(writer, err.Error(), statusCode)
			return
		}

		writer.WriteHeader(statusCode)
		writer.Write(response)
	case http.MethodDelete:
		response, statusCode, err := recordRouter.recordService.AddPrivilegedUser(owner, key, user)
		if err != nil {
			http.Error(writer, err.Error(), statusCode)
			return
		}

		writer.WriteHeader(statusCode)
		writer.Write(response)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("Method not allowed"))
	}
}

func Router(tracker *sync.WaitGroup, logger *log.Logger, db *DBCLient) *http.ServeMux {
	router := http.NewServeMux()

	recordRouter := RecordRouter{
		waitgroup:     tracker,
		recordService: &RecordService{dbClient: db, logger: logger}}

	router.HandleFunc("/record", recordRouter.ServeRecordHTTP)

	router.HandleFunc("/record/allowedUser", recordRouter.ServePrivilegedUserHTTP)

	return router
}
