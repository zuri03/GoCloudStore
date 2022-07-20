package records

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type ParamsMiddleware func(username, password string) error

type Response struct {
	Id string `json:"id"`
}

type AuthHandler struct {
	dbClient       Mongo
	routineTracker *sync.WaitGroup
	validateParams ParamsMiddleware
}

func (handler *AuthHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("IN AUTH HANDLER")
	handler.routineTracker.Add(1)
	defer handler.routineTracker.Done()
	username := req.FormValue("username")
	password := req.FormValue("password")

	if err := handler.validateParams(username, password); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	handler.Authenticate(username, password, writer)
}

func (handler *AuthHandler) Authenticate(username, password string, writer http.ResponseWriter) {
	potentialUsers, err := handler.dbClient.SearchUser(username, password)
	if err != nil {
		fmt.Printf("Error on user search: %s\n", err.Error())
		http.Error(writer, "Internal Server Error on database user search", http.StatusInternalServerError)
		return
	}

	for _, potentialMatch := range potentialUsers {
		matchPassword := potentialMatch.Password
		err := bcrypt.CompareHashAndPassword(matchPassword, []byte(password))
		if err == nil {
			response := Response{Id: potentialMatch.Id}
			jsonBytes, err := json.Marshal(response)
			if err != nil {
				http.Error(writer, "Internal Server Error: error creating json", http.StatusInternalServerError)
				return
			}
			writer.Write(jsonBytes)
			return
		}
	}

	jsonBytes, err := json.Marshal(Response{Id: ""})
	writer.WriteHeader(http.StatusNotFound)
	writer.Write(jsonBytes)
}
