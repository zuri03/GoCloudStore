package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/zuri03/GoCloudStore/common"
)

type FileServerclient struct {
	Client http.Client
}

func NewFileServerClient() *FileServerclient {
	return nil
}

func (client *FileServerclient) SendFile(filePath string) error {
	return nil
}

/*
*Retrieves a file from a server with ID fileId
	- fileId: unique identifier of file that needs to be retrieved from server

Returns:
	- string: file path pointing to tmp file downloaded from the server and saved on this machine
*/
func (client *FileServerclient) GetFile(fileId string) (string, error) {
	return "", nil
}

type RecordServerclient struct {
	HttpClient http.Client
}

func NewRecordServerClient() *RecordServerclient {

	recordServerClient := RecordServerclient{
		HttpClient: http.Client{
			Timeout: time.Duration(time.Second * 60),
		},
	}

	return &recordServerClient
}

type AuthenticationResponse struct {
	Id string `json:"id"`
}

func (recordServerClient *RecordServerclient) AuthenticateUser(username string, password string) (string, bool, error) {
	url := fmt.Sprintf("http://localhost:8080/auth?username=%s&password=%s", username, password)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", false, err
	}

	var authResponse AuthenticationResponse
	if err := recordServerClient.sendRequestAndParseJson(request, authResponse); err != nil {
		return "", false, err
	}

	return authResponse.Id, authResponse.Id != "", nil
}

//returns ID of newly created user
func (recordServerClient *RecordServerclient) CreateUser(username string, password string) (string, error) {
	url := fmt.Sprintf("http://localhost:8080/user?username=%s&password=%s", username, password)
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	var user common.User
	if err := recordServerClient.sendRequestAndParseJson(request, user); err != nil {
		return "", err
	}

	return user.Id, nil
}

func (recordServerClient *RecordServerclient) GetFileRecord(owner string, key string) (*common.Record, error) {
	url := fmt.Sprintf("http://localhost:8080/record?id=%s&key=%s", owner, key)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var record common.Record
	err = recordServerClient.sendRequestAndParseJson(request, record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

//returns error in case file record was unable to be deleted
func (recordServerClient *RecordServerclient) DeleteFileRecord(owner string, key string) error {
	url := fmt.Sprintf("http://localhost:8080/record?owner=%s&key=%s", owner, key)
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	return recordServerClient.sendHttpRequest(request)
}

type FileRecordRequest struct {
	Owner    string `json:"owner"`
	Key      string `json:"key"`
	FileName string `json:"name"`
	Size     int    `json:"size"`
}

func (recordServerClient *RecordServerclient) CreateFileRecord(owner, key, fileName string, fileSize int64) error {
	record := FileRecordRequest{
		Owner:    owner,
		Key:      key,
		FileName: fileName,
		Size:     int(fileSize),
	}
	//need to include necessary params to pass params check middleware
	url := fmt.Sprintf("http://localhost:8080/record?owner=%s&key=%s", owner, key)
	recordBytes, _ := json.Marshal(record)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(recordBytes))
	if err != nil {
		return err
	}

	return recordServerClient.sendHttpRequest(request)
}

func (recordServerClient *RecordServerclient) AddAllowedUser(owner, key, allowedUser string) error {
	url := fmt.Sprintf("http://localhost:8080/record/allowedUser?allowedUser=%s&owner=%s&key=%s", allowedUser, owner, key)

	request, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}

	return recordServerClient.sendHttpRequest(request)
}

func (recordServerClient *RecordServerclient) RemoveAllowedUser(owner, key string, removedUser string) error {

	url := fmt.Sprintf("http://localhost:8080/record/allowedUser?removedUser=%s&owner=%s&key=%s", removedUser, owner, key)

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	return recordServerClient.sendHttpRequest(request)
}

func (recordServerClient *RecordServerclient) sendHttpRequest(request *http.Request) error {
	response, err := recordServerClient.HttpClient.Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		errorMessage, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%d: %s\n", response.StatusCode, string(errorMessage))
	}

	return nil
}

func (recordServerClient *RecordServerclient) sendRequestAndParseJson(request *http.Request, object interface{}) error {
	response, err := recordServerClient.HttpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errorMessage, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%d: %s\n", response.StatusCode, string(errorMessage))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &object); err != nil {
		return err
	}

	return nil
}
