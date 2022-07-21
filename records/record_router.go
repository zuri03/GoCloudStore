package records

import (
	"net/http"
	"sync"

	"github.com/zuri03/GoCloudStore/common"
)

type userDataBase interface {
	GetUser(id string) (*common.User, error)
	CreateUser(user *common.User) error
	SearchUser(username, password string) ([]*common.User, error)
}

type recordDataBase interface {
	GetRecord(key string) (*common.Record, error)
	CreateRecord(record common.Record) error
	DeleteRecord(key string) error
	AddAllowedUser(key, user string) error
	RemoveAllowedUser(key, user string) error
}

func Router(userDB userDataBase, recordDB recordDataBase, tracker *sync.WaitGroup) *http.ServeMux {
	authHandler := AuthHandler{
		dbClient:       userDB,
		routineTracker: tracker,
		validateParams: checkParamsUsername,
	}
	userHandler := UserHandler{
		dbClient:         userDB,
		routineTracker:   tracker,
		paramsMiddleware: checkParamsUsername,
		idMiddleware:     validateId,
	}
	allowedUserHanlder := AllowedUserHandler{
		dbClient:           recordDB,
		routineTracker:     tracker,
		paramsMiddleware:   checkParamsAllowedUser,
		resourceMiddleware: recordExists,
		ownerMiddleware:    checkOwner,
	}
	recordHanlder := RecordHandler{
		dbClient:                 recordDB,
		routineTracker:           tracker,
		paramsMiddleware:         checkParamsRecords,
		resourseExistsMiddleware: recordExists,
		canViewMiddleware:        canView,
		checkOwnerMiddleware:     checkOwner,
	}

	router := http.NewServeMux()

	router.HandleFunc("/record", recordHanlder.ServeHTTP)

	router.HandleFunc("/record/allowedUser", allowedUserHanlder.ServeHTTP)

	router.HandleFunc("/user", userHandler.ServeHTTP)

	router.HandleFunc("/auth", authHandler.ServeHTTP)

	return router
}
