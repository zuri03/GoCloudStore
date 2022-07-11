package records

import (
	"fmt"
	"net/http"
	"strings"
)

/*
The middleware performs a series of checks that ensure the request is in a specific state within each handler
This is why the handlers will return a 500 error when calling keeper.get for example because an "unauthorized" or "not found" error should have been caught in the middleware pipeline
*/

//Checks if a user with the given credentials exists within the system
func authenticate(u *Users, writer http.ResponseWriter, req *http.Request) bool {

	username := req.FormValue("username")
	password := req.FormValue("password")

	if exists, err := u.exists(username, password); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return false
	} else if !exists {
		writer.WriteHeader(http.StatusUnauthorized)
		return false
	}

	return true
}

//Checks if a user can view a given record this is for get requests
func canView(u *Users, r *RecordKeeper, writer http.ResponseWriter, req *http.Request) bool {

	username := req.FormValue("username")
	password := req.FormValue("password")
	key := req.FormValue("key")

	user, err := u.get(username, password)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	authorized, err := r.Authorized(key, user.Id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if !authorized {
		writer.WriteHeader(http.StatusForbidden)
		return false
	}

	return true
}

//Checks if the user is the owner of the record
func checkOwner(u *Users, r *RecordKeeper, writer http.ResponseWriter, req *http.Request) bool {

	username := req.FormValue("username")
	password := req.FormValue("password")
	key := req.FormValue("key")

	user, err := u.get(username, password)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	isOnwer, err := r.IsOnwer(key, user.Id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if !isOnwer {
		writer.WriteHeader(http.StatusForbidden)
		return false
	}

	return true
}

func resourceExists(r *RecordKeeper, writer http.ResponseWriter, req *http.Request) bool {
	key := req.FormValue("key")
	if !r.Exists(key) {
		writer.WriteHeader(http.StatusNotFound)
		return false
	}

	return true
}

//Checks the request for all of the required params
func checkParamsRecords(writer http.ResponseWriter, req *http.Request) bool {
	username := req.FormValue("username")
	password := req.FormValue("password")
	key := req.FormValue("key")

	if key == "" || password == "" || username == "" {
		missing := []string{}
		if key == "" {
			missing = append(missing, "key")
		}
		if password == "" {
			missing = append(missing, "password")
		}
		if username == "" {
			missing = append(missing, "username")
		}
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf("%s missing from request", strings.Join(missing, ","))))

		return false
	}
	return true
}

//Checks request to the user endpoint for the required params
func checkParamsUsers(writer http.ResponseWriter, req *http.Request) bool {
	username := req.FormValue("username")
	password := req.FormValue("password")

	if password == "" || username == "" {
		missing := []string{}
		if password == "" {
			missing = append(missing, "password")
		}
		if username == "" {
			missing = append(missing, "username")
		}
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf("%s missing from request", strings.Join(missing, ","))))

		return false
	}
	return true
}
