package records

import (
	"net/http"
)

type CreateReqest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Key      string `json:"key"`
	FileName string `json:"name"`
	Size     int    `json:"size"`
}

func InitServer(keeper *RecordKeeper) {
	getHandler := GetHandler{Keeper: keeper}
	createHandler := PostHandler{Keeper: keeper}
	router := http.NewServeMux()

	router.HandleFunc("/file", func(writer http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			createHandler.ServeHTTP(writer, req)
		case http.MethodGet:
			getHandler.ServeHTTP(writer, req)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte("Method not allowed"))
		}
	})

}
