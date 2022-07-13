package records

import (
	"fmt"
	"net/http"

	"github.com/zuri03/GoCloudStore/records/db"
)

func Router(keeper *RecordKeeper, mongo *db.Mongo) *http.ServeMux {
	fmt.Println("CREATING SERVER")

	authHandler := AuthHandler{users: mongo}
	userHandler := UserHandler{users: mongo}
	allowedUserHanlder := AllowedUserHandler{Keeper: keeper, Users: mongo}
	recordHanlder := RecordHandler{Keeper: keeper, Users: mongo}
	router := http.NewServeMux()

	router.HandleFunc("/record", recordHanlder.ServeHTTP)

	router.HandleFunc("/record/allowedUser", allowedUserHanlder.ServeHTTP)

	router.HandleFunc("/user", userHandler.ServeHTTP)

	router.HandleFunc("/auth", authHandler.ServeHTTP)

	return router
}
