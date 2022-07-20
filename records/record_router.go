package records

import (
	"net/http"
	"sync"

	"github.com/zuri03/GoCloudStore/common"
)

type Mongo interface {
	GetUser(id string) (*common.User, error)
	CreateUser(user *common.User) error
	SearchUser(username, password string) ([]*common.User, error)

	GetRecord(key string) (*common.Record, error)
	CreateRecord(record common.Record) error
	DeleteRecord(key string) error
	AddAllowedUser(key, user string) error
	RemoveAllowedUser(key, user string) error
}

func Router(mongo Mongo, tracker *sync.WaitGroup) *http.ServeMux {

	authHandler := AuthHandler{dbClient: mongo, routineTracker: tracker, validateParams: checkParamsUsername}
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
