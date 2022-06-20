package records

import (
	"fmt"
	"log"
	"net/http"
)

func InitServer(keeper *RecordKeeper) {
	fmt.Println("CREATING SERVER")
	getHandler := GetHandler{Keeper: keeper}
	createHandler := PostHandler{Keeper: keeper}
	deleteHandler := DeleteHandler{Keeper: keeper}
	addUserHandler := AddUserHandler{Keeper: keeper}
	removeUserHandler := RemoveUserHandler{Keeper: keeper}
	router := http.NewServeMux()

	router.HandleFunc("/record", func(writer http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			createHandler.ServeHTTP(writer, req)
		case http.MethodGet:
			getHandler.ServeHTTP(writer, req)
		case http.MethodDelete:
			deleteHandler.ServeHTTP(writer, req)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte("Method not allowed"))
		}
	})

	router.HandleFunc("/record/allowedUsers", func(writer http.ResponseWriter, req *http.Request) {
		fmt.Println("IN ROUTER")
		switch req.Method {
		case http.MethodPut:
			addUserHandler.ServeHTTP(writer, req)
		case http.MethodDelete:
			removeUserHandler.ServeHTTP(writer, req)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte("Method not allowed"))
		}
	})

	fmt.Println("STARTING LISTENER")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Error occured initializing server: %s\n", err.Error())
		return
	}
	fmt.Printf("Listening on port 8080 \n")

	fmt.Println("INIT SERVER RETURNING")
}
