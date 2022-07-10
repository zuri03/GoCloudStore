package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/zuri03/GoCloudStore/records"
)

func main() {
	fmt.Println("CREATING META DATA SERVER")
	keeper := records.InitRecordKeeper()
	fmt.Println("CREATED META DATA SERVER")
	users := records.NewUsers()
	router := records.Router(&keeper, users)

	server := &http.Server{
		Addr:        ":8080",
		Handler:     router,
		IdleTimeout: 60 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	signaler := make(chan os.Signal)
	signal.Notify(signaler, os.Interrupt)
	signal.Notify(signaler, os.Kill)

	<-signaler
	fmt.Println("SHUT DOWN")
}
