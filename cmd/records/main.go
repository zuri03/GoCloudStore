package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	//"github.com/zuri03/GoCloudStore/common"
	"github.com/zuri03/GoCloudStore/records"
)

//Port these over to an .env file
const PORT = 8080
const HOST = ""
const mongoURIKey = "MONGO_URI"

func main() {

	logOutput, err := os.OpenFile("record-log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer logOutput.Close()

	logger := log.New(logOutput, "", log.LstdFlags)

	logger.Println("Connecting to mongodb..")

	uri := os.Getenv(mongoURIKey)
	db, err := records.NewDBClient(uri)
	if err != nil {
		fmt.Printf("Error in main: %s\n", err.Error())
		return
	}

	tracker := new(sync.WaitGroup)
	router := records.Router(tracker, logger, db)

	address := fmt.Sprintf("%s:%d", HOST, PORT)
	server := &http.Server{
		Addr:        address,
		Handler:     router,
		IdleTimeout: 60 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	logger.Printf("Record server listening on port %d\n", PORT)
	fmt.Printf("Record server listening on port %d\n", PORT)
	signaler := make(chan os.Signal)
	signal.Notify(signaler, os.Interrupt)
	signal.Notify(signaler, os.Kill)

	<-signaler

	logger.Println("Shutdown signal received, waiting for go routines to finish...")

	tracker.Wait()

	logger.Println("Exiting...")
}
