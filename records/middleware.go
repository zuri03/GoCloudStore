package records

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

/*
The middleware performs a series of checks that ensure the request is in a specific state within each handler
This is why the handlers will return a 500 error when calling keeper.get for example because an "unauthorized" or "not found" error should have been caught in the middleware pipeline
*/

//Checks if a user with the given credentials exists within the system
/*
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
*/

//Checks if a user can view a given record this is for get requests
func canView(id string, key string, r *RecordKeeper, writer http.ResponseWriter) bool {
	authorized, err := r.Authorized(key, id)
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
func checkOwner(id, key string, r *RecordKeeper, writer http.ResponseWriter) bool {
	isOnwer, err := r.IsOnwer(key, id)
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

func resourceExists(key string, r *RecordKeeper, writer http.ResponseWriter) bool {
	if !r.Exists(key) {
		writer.WriteHeader(http.StatusNotFound)
		return false
	}

	return true
}

//Checks the request for all of the required params
func checkParamsRecords(writer http.ResponseWriter, req *http.Request) bool {
	id := req.FormValue("owner")
	key := req.FormValue("key")

	if key == "" {
		missing := []string{}
		if key == "" {
			missing = append(missing, "key")
		}
		if id == "" {
			missing = append(missing, "id")
		}
		http.Error(writer, fmt.Sprintf("%s missing from request", strings.Join(missing, ",")), http.StatusBadRequest)

		return false
	}

	return validateId(id, writer)
}

//Checks request to the user endpoint for the required params
func checkParamsUsername(writer http.ResponseWriter, req *http.Request) bool {
	username := req.FormValue("username")
	password := req.FormValue("password")

	fmt.Println("CHECKING USERNAME/PASS")
	fmt.Printf("Username => %s | Pass => %s\n", username, password)

	if password == "" || username == "" {
		missing := []string{}
		if password == "" {
			missing = append(missing, "password")
		}
		if username == "" {
			missing = append(missing, "username")
		}
		http.Error(writer, fmt.Sprintf("%s missing from request", strings.Join(missing, ",")), http.StatusBadRequest)
		return false
	}
	return true
}

func checkParamsId(writer http.ResponseWriter, req *http.Request) bool {
	id := req.FormValue("id")

	if id == "" {
		http.Error(writer, "Error: Id is mssing", http.StatusBadRequest)
		return false
	}
	return validateId(id, writer)
}

func checkParamsAllowedUser(writer http.ResponseWriter, req *http.Request) bool {
	user := req.FormValue("user")

	if user == "" {
		http.Error(writer, "Missing user parameter", http.StatusBadRequest)
		return false
	}

	return validateId(user, writer)
}

func validateId(id string, writer http.ResponseWriter) bool {
	if id == "" {
		http.Error(writer, "Request is missing id", http.StatusBadRequest)
		return false
	}

	if _, err := uuid.Parse(id); err != nil {
		http.Error(writer, "Unable to parse id", http.StatusBadRequest)
		return false
	}

	return true
}
