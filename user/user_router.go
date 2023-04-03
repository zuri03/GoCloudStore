package user

import (
	"fmt"
	"net/http"

	"github.com/zuri03/GoCloudStore/common"
)

type userDataBase interface {
	GetUser(id string) (*common.User, error)
	CreateUser(user *common.User) error
	SearchUser(username, password string) ([]*common.User, error)
}

func ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	//handler.routineTracker.Add(1)
	//defer handler.routineTracker.Done()

	fmt.Printf("req => %s\n", req.Method)

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

func Router() *http.ServeMux {

	router := http.NewServeMux()

	router.HandleFunc("/user", ServeHTTP)

	return router
}
