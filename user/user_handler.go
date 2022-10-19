package user

import (
	"fmt"
	"net/http"
	"sync"
)

type userHandler struct {
	db               userDataBase
	routineTracker   *sync.WaitGroup
	paramsMiddleware func(username, password string) error
	idMiddleware     func(id string) error
}

func (handler *userHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	//handler.routineTracker.Add(1)
	//defer handler.routineTracker.Done()

	switch req.Method {
	case http.MethodPost:
		writer.Write([]byte("POST METHOD"))
	case http.MethodGet:
		writer.Write([]byte("POST GET"))
	case http.MethodDelete:
		writer.Write([]byte("POST DELETE"))
	default:
		http.Error(writer, fmt.Sprintf("Method %s is not supported on this endpoint", req.Method), http.StatusMethodNotAllowed)
	}
}
