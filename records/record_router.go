package records

import (
	"net/http"
	"sync"

	"github.com/zuri03/GoCloudStore/records/db"
)

func Router(keeper *RecordKeeper, mongo *db.Mongo, tracker *sync.WaitGroup) *http.ServeMux {

	authHandler := AuthHandler{users: mongo, routineTracker: tracker}
	userHandler := UserHandler{users: mongo, routineTracker: tracker}
	allowedUserHanlder := AllowedUserHandler{
		Keeper: keeper, Users: mongo, routineTracker: tracker}
	recordHanlder := RecordHandler{
		Keeper: keeper, Users: mongo, routineTracker: tracker}

	router := http.NewServeMux()

	router.HandleFunc("/record", recordHanlder.ServeHTTP)

	router.HandleFunc("/record/allowedUser", allowedUserHanlder.ServeHTTP)

	router.HandleFunc("/user", userHandler.ServeHTTP)

	router.HandleFunc("/auth", authHandler.ServeHTTP)

	return router
}
