package user

import (
	"net/http"

	"github.com/zuri03/GoCloudStore/common"
)

type userDataBase interface {
	GetUser(id string) (*common.User, error)
	CreateUser(user *common.User) error
	SearchUser(username, password string) ([]*common.User, error)
}

func Router() *http.ServeMux {

	router := http.NewServeMux()

	userHandler := userHandler{}

	router.HandleFunc("/user", userHandler.ServeHTTP)

	return router
}
