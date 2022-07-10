package records

import (
	"fmt"
	"net/http"
)

func Router(keeper *RecordKeeper, users *Users) *http.ServeMux {
	fmt.Println("CREATING SERVER")
	getHandler := GetHandler{Keeper: keeper, Users: users}
	createHandler := PostHandler{Keeper: keeper, Users: users}
	deleteHandler := DeleteHandler{Keeper: keeper, Users: users}
	addUserHandler := AddUserHandler{Keeper: keeper, Users: users}
	removeUserHandler := RemoveUserHandler{Keeper: keeper, Users: users}
	createUserHandler := CreateHandler{Users: users}
	getUserHandler := GetUserHandler{Users: users}
	authHandler := AuthHandler{Users: users}
	router := http.NewServeMux()

	router.HandleFunc("/record", func(writer http.ResponseWriter, req *http.Request) {
		if hasParams := checkParamsRecords(writer, req); !hasParams {
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

	router.HandleFunc("/record/allowedUser", func(writer http.ResponseWriter, req *http.Request) {
		if hasParams := checkParamsRecords(writer, req); !hasParams {
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

	router.HandleFunc("/user", func(writer http.ResponseWriter, req *http.Request) {
		if hasParams := checkParamsUsers(writer, req); !hasParams {
			return
		}
		switch req.Method {
		case http.MethodPost:
			createUserHandler.ServeHTTP(writer, req)
		case http.MethodGet:
			getUserHandler.ServeHTTP(writer, req)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte("Method not allowed"))
		}
	})

	router.HandleFunc("/auth", func(writer http.ResponseWriter, req *http.Request) {
		authHandler.ServeHTTP(writer, req)
	})

	return router
}
