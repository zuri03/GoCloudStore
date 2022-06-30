package users

import (
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Id           []byte
	Username     string
	Password     string
	CreationDate string
}

const bcryptCost = 10

func InitServer() {
	fmt.Println("CREATING SERVER")
	router := http.NewServeMux()

	users := make(map[string]User)
	createHandler := CreateHandler{Users: users}
	GetHandler := GetHandler{Users: users}

	router.HandleFunc("/user", func(writer http.ResponseWriter, req *http.Request) {
		fmt.Println("GOT USER REQUEST")
		switch req.Method {
		case http.MethodGet:
			GetHandler.ServeHTTP(writer, req)
		case http.MethodPost:
			createHandler.ServeHTTP(writer, req)
		case http.MethodDelete:
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte("Method not allowed"))
		}
	})

	fmt.Println("STARTING LISTENER")
	err := http.ListenAndServe(":9000", router)
	if err != nil {
		log.Fatalf("Error occured initializing server: %s\n", err.Error())
		return
	}
	fmt.Printf("Listening on port 9000 \n")

	fmt.Println("INIT SERVER RETURNING")
}
