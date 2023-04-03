package user

import (
	"fmt"
	"net/http"
	"sync"
)

type UserRouter struct {
	waitgroup *sync.WaitGroup
}

func (router *UserRouter) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	router.waitgroup.Add(1)
	defer router.waitgroup.Done()

	switch req.Method {
	case http.MethodPost:
		username := req.FormValue("username")
		password := req.FormValue("password")

		newUser, statusCode, err := CreateUser(username, password)
		if err != nil {
			http.Error(writer, err.Error(), statusCode)
			return
		}

		writer.WriteHeader(statusCode)
		writer.Write(newUser)
		return
	case http.MethodGet:
		id := req.FormValue("id")

		user, statusCode, err := GetUser(id)
		if err != nil {
			http.Error(writer, err.Error(), statusCode)
			return
		}

		writer.WriteHeader(statusCode)
		writer.Write(user)
		return
	case http.MethodDelete:
		id := req.FormValue("id")

		statusCode, err := DeleteUser(id)
		if err != nil {
			http.Error(writer, err.Error(), statusCode)
			return
		}

		writer.WriteHeader(statusCode)
		return
	default:
		http.Error(writer, fmt.Sprintf("Method %s is not supported on this endpoint", req.Method), http.StatusMethodNotAllowed)
	}
}

func Router(waitgroup *sync.WaitGroup) *http.ServeMux {
	router := http.NewServeMux()

	userRouter := UserRouter{waitgroup: waitgroup}

	router.HandleFunc("/user", userRouter.ServeHTTP)

	return router
}
