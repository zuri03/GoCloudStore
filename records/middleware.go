package records

import (
	"net/http"
)

func authenticate(u *UserClient, writer http.ResponseWriter, req *http.Request) (bool, error) {

	username := req.FormValue("username")
	password := req.FormValue("password")

	if exists, err := u.Exists(username, password); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return false, err
	} else if !exists {
		writer.WriteHeader(http.StatusUnauthorized)
		return false, nil
	}

	return true, nil
}

func checkParams(writer http.ResponseWriter, req *http.Request) bool {
	username := req.FormValue("username")
	password := req.FormValue("password")
	key := req.FormValue("key")

	if key == "" || password == "" || username == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Key, Username or password missing from request"))
		return false
	}
	return true
}
