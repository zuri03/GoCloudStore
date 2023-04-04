package user

import (
	"fmt"
	"net/http"
	"sync"
)

type UserRouter struct {
	waitgroup *sync.WaitGroup
	service   *UserService
}

func (router *UserRouter) ServeUserHTTP(writer http.ResponseWriter, req *http.Request) {
	router.waitgroup.Add(1)
	defer router.waitgroup.Done()

	switch req.Method {
	case http.MethodPost:
		username := req.FormValue("username")
		password := req.FormValue("password")

		newUser, statusCode, err := router.service.CreateUser(username, password)
		if err != nil {
			http.Error(writer, err.Error(), statusCode)
			return
		}

		writer.WriteHeader(statusCode)
		writer.Write(newUser)
		return
	case http.MethodGet:
		id := req.FormValue("id")

		user, statusCode, err := router.service.GetUser(id)
		if err != nil {
			http.Error(writer, err.Error(), statusCode)
			return
		}

		writer.WriteHeader(statusCode)
		writer.Write(user)
		return
	case http.MethodDelete:
		id := req.FormValue("id")

		statusCode, err := router.service.DeleteUser(id)
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

func (router *UserRouter) ServeAuthHTTP(writer http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(writer, fmt.Sprintf("Method %s is not supported on this endpoint", req.Method), http.StatusMethodNotAllowed)
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	authResponseBytes, statusCode, err := router.service.AuthorizeUser(username, password)
	if err != nil {
		http.Error(writer, err.Error(), statusCode)
		return
	}

	writer.WriteHeader(statusCode)
	writer.Write(authResponseBytes)
}

func Router(waitgroup *sync.WaitGroup, db *DBClient) *http.ServeMux {
	router := http.NewServeMux()

	userRouter := UserRouter{
		waitgroup: waitgroup,
		service:   &UserService{db: db}}

	router.HandleFunc("/user", userRouter.ServeUserHTTP)

	router.HandleFunc("/login", userRouter.ServeAuthHTTP)

	return router
}
