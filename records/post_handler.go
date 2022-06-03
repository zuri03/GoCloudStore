package records

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PostHandler struct {
	Keeper *RecordKeeper
}

func (handler *PostHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	fmt.Printf("Body => %s\n", string(body))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("INTERNAL SERVER ERROR"))
		return
	}

	var request CreateReqest
	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Printf("Error decoding json => %s\n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("INTERNAL SERVER ERROR"))
		return
	}
	fmt.Printf("DECODED JSON => %s\n", request.Username)

	record, err := handler.Keeper.SetRecord(request.Key, fmt.Sprintf("%s:%s", request.Username, request.Password),
		request.Size, request.FileName)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf("record %s alread exists", request.Key)))
		return
	}
	jsonBytes, _ := json.Marshal(*record)
	writer.Write(jsonBytes)
}
