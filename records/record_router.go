package records

import (
	"fmt"
	"net/http"
)

func Router(keeper *RecordKeeper, users *Users) *http.ServeMux {
	fmt.Println("CREATING SERVER")
	getHandler := GetHandler{Keeper: keeper}
	createHandler := PostHandler{Keeper: keeper}
	deleteHandler := DeleteHandler{Keeper: keeper}
	addUserHandler := AddUserHandler{Keeper: keeper}
	removeUserHandler := RemoveUserHandler{Keeper: keeper}
	router := http.NewServeMux()

	router.HandleFunc("/record", func(writer http.ResponseWriter, req *http.Request) {
		if hasParams := checkParams(writer, req); !hasParams {
			return
		}

		if authenticated, err := authenticate(users, writer, req); !authenticated || err != nil {
			return
		}

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
		if hasParams := checkParams(writer, req); !hasParams {
			return
		}

		if authenticated, err := authenticate(users, writer, req); !authenticated || err != nil {
			return
		}

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

	return router
}
