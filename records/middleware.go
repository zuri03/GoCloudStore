package records

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/zuri03/GoCloudStore/common"
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
func canView(id string, record common.Record, writer http.ResponseWriter) bool {
	fmt.Printf("Checking if %s is owner %s\n", id, record.Owner)
	if record.Owner == id {
		fmt.Println("Found match in owner")
		return true
	}

	for _, allowedUser := range record.AllowedUsers {
		fmt.Printf("%s == %s\n", id, allowedUser)
		if id == allowedUser {
			return true
		}
	}

	return false
}

//Checks if the user is the owner of the record
func checkOwner(id string, record common.Record, writer http.ResponseWriter) bool {
	fmt.Println("CHECKING OWNER")
	if id != record.Owner {
		fmt.Printf("Owner %s does not match user %s\n", record.Owner, id)
		http.Error(writer, "User is not authorized", http.StatusUnauthorized)
		return false
	}
	fmt.Println("Owner is user match")
	return true
}

func resourceExists(key string, db Mongo, writer http.ResponseWriter) (common.Record, bool) {
	record, err := db.GetRecord(key)
	//If the method returns an error then something serious has occured
	//If the method cannot find the resource it returns an empty struct with a nil error
	if err != nil {
		fmt.Printf("Error in checking resource: %s\n", err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return common.Record{}, false
	}

	if record.Key == "" {
		message := fmt.Sprintf("Unable to find record %s", key)
		http.Error(writer, message, http.StatusNotFound)
		return common.Record{}, false
	}

	return *record, true
}

//Checks the request for all of the required params
func checkParamsRecords(writer http.ResponseWriter, req *http.Request) bool {
	id := req.FormValue("id")
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
	owner := req.FormValue("owner")
	key := req.FormValue("key")

	if user == "" || owner == "" || key == "" {
		missing := []string{}
		if user == "" {
			missing = append(missing, "user")
		}
		if owner == "" {
			missing = append(missing, "owner")
		}
		if key == "" {
			missing = append(missing, "key")
		}
		http.Error(writer, fmt.Sprintf("%s missing from request", strings.Join(missing, ",")), http.StatusBadRequest)
		return false
	}

	return validateId(user, writer) && validateId(owner, writer)
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
