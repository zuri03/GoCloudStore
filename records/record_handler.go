package records

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/zuri03/GoCloudStore/common"
)

type RecordHandler struct {
	dbClient       Mongo
	routineTracker *sync.WaitGroup
}

type Request struct {
	Owner    string `json:"owner"`
	Key      string `json:"key"`
	FileName string `json:"name"`
	Size     int    `json:"size"`
}

//refactor this handler
func (handler *RecordHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	handler.routineTracker.Add(1)
	defer handler.routineTracker.Done()

	if req.Method == http.MethodPost {
		var requestBody Request
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(writer, "Unable to read request body", http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(body, &requestBody); err != nil {
			http.Error(writer, "Unable to read request body", http.StatusBadRequest)
			return
		}

		handler.CreateRecord(requestBody, writer)
		return
	}

	if !checkParamsRecords(writer, req) {
		return
	}

	//autheticate the user here
	id := req.FormValue("id")
	key := req.FormValue("key")

	record, ok := resourceExists(key, handler.dbClient, writer)
	if !ok {
		return
	}

	switch req.Method {
	case http.MethodGet:
		if !canView(id, record, writer) {
			return
		}
		handler.GetRecord(record, writer)
	case http.MethodDelete:
		if !checkOwner(id, record, writer) {
			return
		}
		handler.DeleteRecord(key, writer)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("Method not allowed"))
	}
}

func (handler *RecordHandler) CreateRecord(request Request, writer http.ResponseWriter) {
	location := fmt.Sprintf("%s/%s", request.Owner, request.FileName)
	currentTime := time.Now().Format("2006-01-02 03:04:05")

	newRecord := common.Record{
		Size:         int64(request.Size),
		Name:         request.FileName,
		Location:     location,
		CreatedAt:    currentTime,
		IsPublic:     false,
		Owner:        request.Owner,
		AllowedUsers: make([]string, 0),
	}

	if err := handler.dbClient.CreateRecord(newRecord); err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(newRecord)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Write(jsonBytes)
}

func (handler *RecordHandler) GetRecord(record common.Record, writer http.ResponseWriter) {
	jsonBytes, err := json.Marshal(record)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Write(jsonBytes)
}

func (handler *RecordHandler) DeleteRecord(key string, writer http.ResponseWriter) {
	if err := handler.dbClient.DeleteRecord(key); err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
