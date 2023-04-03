package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/zuri03/GoCloudStore/user"
)

func main() {

	waitgroup := new(sync.WaitGroup)
	router := user.Router(waitgroup)

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

	fmt.Println("Shutdown signal received, waiting for go routines to finish...")

	waitgroup.Wait()

	fmt.Println("Exiting...")
}
