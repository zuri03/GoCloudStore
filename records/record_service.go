package records

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zuri03/GoCloudStore/common"
)

type Request struct {
	Owner    string `json:"owner"`
	Key      string `json:"key"`
	FileName string `json:"name"`
	Size     int    `json:"size"`
}

type RecordService struct {
	dbClient *DBCLient
	logger   *log.Logger
}

func (recordService *RecordService) GetRecord(id string, key string) ([]byte, int, error) {
	if id == "" || key == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("Error: id or key missing from request got, id: %s, key %s\n", id, key)
	}

	if _, err := uuid.Parse(id); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("id is not valid uuid got, %s\n", id)
	}

	record, err := recordService.dbClient.GetRecord(key)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	//the db client could not find a record so it returned an empty struct
	if record.Key == "" {
		return nil, http.StatusNotFound, fmt.Errorf("Error: could not find record with key %s\n", key)
	}

	if !record.IsPublic && record.Owner != id {
		var isAllowedUser bool = false
		for _, allowedUser := range record.AllowedUsers {
			if id == allowedUser {
				isAllowedUser = true
			}
		}

		if !isAllowedUser {
			return nil, http.StatusForbidden, fmt.Errorf("Error: id %s does not have access to record\n", id)
		}
	}

	jsonBytes, err := json.Marshal(record)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return jsonBytes, http.StatusOK, nil
}

func (recordService *RecordService) DeleteRecord(id string, key string) (int, error) {
	if id == "" || key == "" {
		return http.StatusBadRequest, fmt.Errorf("Error: id or key missing from request got, id: %s, key %s\n", id, key)
	}

	if _, err := uuid.Parse(id); err != nil {
		return http.StatusBadRequest, fmt.Errorf("id is not valid uuid got, %s\n", id)
	}

	record, err := recordService.dbClient.GetRecord(key)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if record.Owner != id {
		return http.StatusForbidden, fmt.Errorf("Error: id %s does not have access to record\n", id)
	}

	if err := recordService.dbClient.DeleteRecord(key); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (recordService *RecordService) CreateRecord(newRecordRequest RecordCreationRequest) ([]byte, int, error) {

	if newRecordRequest.Key == "" || newRecordRequest.Owner == "" || newRecordRequest.FileName == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("Error: owner, key or file name missing from request got, owner: %s, key %s, file name: %s\n", newRecordRequest.Owner, newRecordRequest.Key, newRecordRequest.FileName)
	}

	potentialMatch, err := recordService.dbClient.GetRecord(newRecordRequest.Key)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if potentialMatch != nil || potentialMatch.Key != "" {
		return nil, http.StatusConflict, fmt.Errorf("Error: record with key %s already exists\n", potentialMatch.Key)
	}

	location := fmt.Sprintf("%s/%s", newRecordRequest.Owner, newRecordRequest.FileName)
	currentTime := time.Now().Format("2006-01-02 03:04:05")

	newRecord := common.Record{
		Key:          newRecordRequest.Key,
		Size:         int64(newRecordRequest.Size),
		Name:         newRecordRequest.FileName,
		Location:     location,
		CreatedAt:    currentTime,
		IsPublic:     false,
		Owner:        newRecordRequest.Owner,
		AllowedUsers: make([]string, 0),
	}

	if err := recordService.dbClient.CreateRecord(newRecord); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	jsonBytes, err := json.Marshal(newRecord)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return jsonBytes, http.StatusOK, nil
}
