package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/zuri03/GoCloudStore/records"
	"github.com/zuri03/GoCloudStore/records/db"
)

const PORT = 8080
const HOST = ""

func main() {

	logOutput, err := os.OpenFile("record-log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(logOutput)

	log.Println("Connecting to mongodb..")
	mongo, err := db.New()
	if err != nil {
		log.Fatalf("Error connecting to mongo: %s\n", err.Error())
		return
	}
	log.Println("Successfully connected to mongodb")

	tracker := new(sync.WaitGroup)
	router := records.Router(mongo, mongo, tracker)

	address := fmt.Sprintf("%s:%d", HOST, PORT)
	server := &http.Server{
		Addr:        address,
		Handler:     router,
		IdleTimeout: 60 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
		
	}()
	log.Printf("Record server listening on port %d\n", PORT)
	signaler := make(chan os.Signal)
	signal.Notify(signaler, os.Interrupt)
	signal.Notify(signaler, os.Kill)

	<-signaler

	log.Println("Shutdown signal received, waiting for go routines to finish...")

	tracker.Wait()

	fmt.Println("Exiting...")
}
