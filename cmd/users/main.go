package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/zuri03/GoCloudStore/user"
	//"github.com/zuri03/GoCloudStore/user/db"
)

func main() {

	/*userDb, err := db.New()

	if err != nil {
		fmt.Printf("Error starting db %s", err.Error())
	}
	*/
	router := user.Router()

	fmt.Printf("%t", router == nil)
	server := &http.Server{
		Addr:        ":9000",
		Handler:     router,
		IdleTimeout: 60 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	fmt.Println("SERVER LISTENING ON PORT 9000")

	signaler := make(chan os.Signal)
	signal.Notify(signaler, os.Interrupt)
	signal.Notify(signaler, os.Kill)

	<-signaler

	fmt.Println("Exiting...")
}
