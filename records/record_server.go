package records

import (
	"fmt"
	"log"
	"net/http"
)

type CreateReqest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Key      string `json:"key"`
	FileName string `json:"name"`
	Size     int64  `json:"size"`
}

func InitServer(keeper *RecordKeeper) {
	fmt.Println("CREATING SERVER")
	getHandler := GetHandler{Keeper: keeper}
	createHandler := PostHandler{Keeper: keeper}
	router := http.NewServeMux()

	router.HandleFunc("/file", func(writer http.ResponseWriter, req *http.Request) {
		fmt.Println("IN ROUTER")
		switch req.Method {
		case http.MethodPost:
			fmt.Println("IN POST REQUEST")
			createHandler.ServeHTTP(writer, req)
		case http.MethodGet:
			fmt.Println("IN GET REQUEST")
			getHandler.ServeHTTP(writer, req)
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
