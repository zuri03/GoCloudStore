package records

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CreateReqest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Key      string `json:"key"`
	FileName string `json:"name"`
	Size     int64  `json:"size"`
}

type PostHandler struct {
	Keeper *RecordKeeper
	Users  *Users
}

func (handler *PostHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	fmt.Printf("Body => %s\n", string(body))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var request CreateReqest
	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Printf("Error decoding json => %s\n", err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Incorrect json formatting"))
		return
	}
	owner, err := handler.Users.get(request.Username, request.Password)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	record, err := handler.Keeper.New(request.Key, owner.Id, request.Size, request.FileName)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf("record %s alread exists", request.Key)))
		return
	}
	jsonBytes, _ := json.Marshal(*record)
	writer.Write(jsonBytes)
}
