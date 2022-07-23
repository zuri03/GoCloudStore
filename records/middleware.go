package records

import (
	"fmt"
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
func canView(id string, record common.Record) error {
	if record.IsPublic {
		return nil
	}

	if record.Owner == id {
		fmt.Println("Found match in owner")
		return nil
	}

	for _, allowedUser := range record.AllowedUsers {
		fmt.Printf("%s == %s\n", id, allowedUser)
		if id == allowedUser {
			return nil
		}
	}

	return fmt.Errorf("user is not allowed to view this record")
}

//Checks if the user is the owner of the record
func checkOwner(id string, record common.Record) error {
	if id != record.Owner {
		return fmt.Errorf("User is not authorized")
	}
	return nil
}

func recordExists(key string, db recordDataBase) (common.Record, error) {
	record, err := db.GetRecord(key)
	//If the method returns an error then something serious has occured
	//If the method cannot find the resource it returns an empty struct with a nil error
	if err != nil {
		fmt.Printf("Error in checking resource: %s\n", err.Error())
		return common.Record{}, fmt.Errorf("Internal Server Error")
	}

	if record.Key == "" {
		fmt.Printf("Unable to find record %s", key)
		return common.Record{}, nil
	}

	return *record, nil
}

//Checks the request for all of the required params
func checkParamsRecords(id, key string) error {

	if key == "" {
		missing := []string{}
		if key == "" {
			missing = append(missing, "key")
		}
		if id == "" {
			missing = append(missing, "id")
		}
	}

	return validateId(id)
}

//Checks request to the user endpoint for the required params
func checkParamsUsername(username, password string) error {
	if password != "" && username != "" {
		return nil
	}

	missing := []string{}
	if password == "" {
		missing = append(missing, "password")
	}
	if username == "" {
		missing = append(missing, "username")
	}

	return fmt.Errorf("%s missing from request", strings.Join(missing, ","))
}

func checkParamsAllowedUser(user, owner, key string) error {

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
		return fmt.Errorf("%s missing from request", strings.Join(missing, ","))
	}
	if err := validateId(user); err != nil {
		return err
	}
	if err := validateId(owner); err != nil {
		return err
	}

	return nil
}

func validateId(id string) error {
	if id == "" {
		return fmt.Errorf("id is missing")
	}

	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("id is not valid uuid")
	}

	return nil
}
