package records

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/zuri03/GoCloudStore/common"
)

type Request struct {
	Owner    string `json:"owner"`
	Key      string `json:"key"`
	FileName string `json:"name"`
	Size     int    `json:"size"`
}

type RecordHandler struct {
	dbClient                 recordDataBase
	routineTracker           *sync.WaitGroup
	paramsMiddleware         func(id, key string) error
	resourseExistsMiddleware func(key string, db recordDataBase) (common.Record, error)
	canViewMiddleware        func(id string, record common.Record) error
	checkOwnerMiddleware     func(id string, record common.Record) error
	logger                   *log.Logger
}

//refactor this handler
func (handler *RecordHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	handler.routineTracker.Add(1)
	defer handler.routineTracker.Done()

	handler.logger.Printf("Record request received, method: %s\n", req.Method)

	if req.Method == http.MethodPost {
		
		var requestBody Request
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			handler.logger.Printf("Error reading body: %s\n", err.Error())
			http.Error(writer, "Unable to read request body", http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(body, &requestBody); err != nil {
			handler.logger.Printf("Error unmarshaling json: %s\n", err.Error())
			http.Error(writer, "Unable to read request body", http.StatusBadRequest)
			return
		}

		record, err := handler.resourseExistsMiddleware(requestBody.Key, handler.dbClient)
		if err != nil {
			handler.logger.Printf("Error in resource middleware: %s\n", err.Error())
			http.Error(writer, fmt.Sprintf("Internal Server Error"), http.StatusInternalServerError)
			return
		}

		if record.Key != "" {
			handler.logger.Printf("Record already exists: %v\n", record)
			http.Error(writer, fmt.Sprintf("%s already exists", requestBody.Key), http.StatusConflict)
			return
		}

		handler.logger.Printf("Post Request Successful")
		handler.CreateRecord(requestBody, writer)
		return
	}

	//autheticate the user here
	id := req.FormValue("id")
	key := req.FormValue("key")

	if err := handler.paramsMiddleware(id, key); err != nil {
		handler.logger.Printf("Error in params middleware: %s\n", err.Error())
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := handler.resourseExistsMiddleware(key, handler.dbClient)
	if err != nil {
		handler.logger.Printf("Error in resource middleware: %s\n", err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if record.Key == "" {
		handler.logger.Printf("Unable to find record: %s\n", key)
		http.Error(writer, fmt.Sprintf("Cannot find %s", key), http.StatusNotFound)
		return
	}

	switch req.Method {
	case http.MethodGet:
		if err := handler.canViewMiddleware(id, record); err != nil {
			http.Error(writer, err.Error(), http.StatusForbidden)
			return
		}
		handler.GetRecord(record, writer)
	case http.MethodDelete:
		if err := handler.checkOwnerMiddleware(id, record); err != nil {
			http.Error(writer, err.Error(), http.StatusForbidden)
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
		Key:          request.Key,
		Size:         int64(request.Size),
		Name:         request.FileName,
		Location:     location,
		CreatedAt:    currentTime,
		IsPublic:     false,
		Owner:        request.Owner,
		AllowedUsers: make([]string, 0),
	}

	if err := handler.dbClient.CreateRecord(newRecord); err != nil {
		handler.logger.Printf("Error creating record: %s\n", err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(newRecord)
	if err != nil {
		handler.logger.Printf("Error marshaling json: %s\n", err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Write(jsonBytes)
}

func (handler *RecordHandler) GetRecord(record common.Record, writer http.ResponseWriter) {
	jsonBytes, err := json.Marshal(record)
	if err != nil {
		handler.logger.Printf("Error marshaling json: %s\n", err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Write(jsonBytes)
}

func (handler *RecordHandler) DeleteRecord(key string, writer http.ResponseWriter) {
	if err := handler.dbClient.DeleteRecord(key); err != nil {
		handler.logger.Printf("Error deleting record: %s\n", err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
