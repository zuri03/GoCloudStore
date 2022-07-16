package records

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/zuri03/GoCloudStore/records/db"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Id string `json:"id"`
}

type AuthHandler struct {
	users          *db.Mongo
	routineTracker *sync.WaitGroup
}

func (handler *AuthHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("IN AUTH HANDLER")
	handler.routineTracker.Add(1)
	defer handler.routineTracker.Done()

	if !checkParamsUsername(writer, req) {
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	handler.Authenticate(username, password, writer)
}

func (handler *AuthHandler) Authenticate(username, password string, writer http.ResponseWriter) {
	fmt.Println("In auth method")
	potentialUsers, err := handler.users.SearchUser(username, password)
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
			}
			writer.Write(jsonBytes)
			return
		}
	}

	fmt.Println("Returining not found")
	jsonBytes, err := json.Marshal(Response{Id: ""})
	writer.WriteHeader(http.StatusNotFound)
	writer.Write(jsonBytes)
}
