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

func (Handler *PostHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("INTERNAL SERVER ERROR"))
		return
	}

	var request CreateReqest
	json.Unmarshal(body, &request)
	fmt.Printf("DECODED JSON => %s", request.Username)
	writer.Write([]byte(request.Username))
}
