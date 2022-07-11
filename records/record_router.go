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
		if !checkParamsRecords(writer, req) {
			return
		}

		if !authenticate(users, writer, req) {
			return
		}

		//For create request the pipeline ends here the request has passed all of the checks
		if req.Method == http.MethodPost {
			createHandler.ServeHTTP(writer, req)
			return
		}

		if !resourceExists(keeper, writer, req) {
			return
		}

		switch req.Method {
		case http.MethodGet:
			if !canView(users, keeper, writer, req) {
				return
			}
			getHandler.ServeHTTP(writer, req)
		case http.MethodDelete:
			if !checkOwner(users, keeper, writer, req) {
				return
			}
			deleteHandler.ServeHTTP(writer, req)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte("Method not allowed"))
		}
	})

	router.HandleFunc("/record/allowedUser", func(writer http.ResponseWriter, req *http.Request) {
		//Middleware pipeline
		if !checkParamsRecords(writer, req) {
			return
		}

		if !authenticate(users, writer, req) {
			return
		}

		if !resourceExists(keeper, writer, req) {
			return
		}

		if !checkOwner(users, keeper, writer, req) {
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
