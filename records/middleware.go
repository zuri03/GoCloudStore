package records

import (
	"fmt"
	"net/http"
	"strings"
)

func authenticate(u *Users, writer http.ResponseWriter, req *http.Request) (bool, error) {

	username := req.FormValue("username")
	password := req.FormValue("password")

	if exists, err := u.exists(username, password); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return false, err
	} else if !exists {
		writer.WriteHeader(http.StatusUnauthorized)
		return false, nil
	}

	return true, nil
}

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
