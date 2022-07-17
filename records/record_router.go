package records

import (
	"net/http"
	"sync"
)

type Mongo interface {
	GetUser(id string) (*User, error)
	CreateUser(user *User) error
	SearchUser(username, password string) ([]*User, error)

	GetRecord(key string) (*Record, error)
	CreateRecord(record Record) error
	DeleteRecord(key string) error
	AllowedUser(key, id, user string) error
	RemoveAllowedUser(key, id, user string) error
}

func Router(mongo Mongo, tracker *sync.WaitGroup) *http.ServeMux {

	authHandler := AuthHandler{dbClient: mongo, routineTracker: tracker}
	userHandler := UserHandler{dbClient: mongo, routineTracker: tracker}
	allowedUserHanlder := AllowedUserHandler{dbClient: mongo, routineTracker: tracker}
	recordHanlder := RecordHandler{dbClient: mongo, routineTracker: tracker}

	router := http.NewServeMux()

	router.HandleFunc("/record", recordHanlder.ServeHTTP)

	router.HandleFunc("/record/allowedUser", allowedUserHanlder.ServeHTTP)

	router.HandleFunc("/user", userHandler.ServeHTTP)

	router.HandleFunc("/auth", authHandler.ServeHTTP)

	return router
}
