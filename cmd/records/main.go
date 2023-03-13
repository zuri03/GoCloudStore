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
	"github.com/zuri03/GoCloudStore/records/db"
)

//Port these over to an .env file
const PORT = 8080
const HOST = ""

//A temporary struct that allows me to remove the dependency on a mongodb instance and run the record server in a k8s service independently
/*
type DB struct{}

func (d DB) GetUser(id string) (*common.User, error) {
	return &common.User{
		Id:           "id",
		Username:     "username",
		Password:     []byte("password"),
		CreationDate: "",
	}, nil
}

func (d DB) CreateUser(user *common.User) error {
	return nil
}

func (d DB) SearchUser(username, password string) ([]*common.User, error) {
	return []*common.User{
		{
			Id:           "id",
			Username:     "username",
			Password:     []byte("password"),
			CreationDate: "",
		},
	}, nil
}

func (d DB) GetRecord(key string) (*common.Record, error) {
	return &common.Record{
		Key: key,
	}, nil
}

func (d DB) CreateRecord(record common.Record) error {
	return nil
}

func (d DB) DeleteRecord(key string) error {
	return nil
}

func (d DB) ReplaceRecord(record *common.Record) error {
	return nil
}
*/

func main() {

	logOutput, err := os.OpenFile("record-log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer logOutput.Close()

	logger := log.New(logOutput, "", log.LstdFlags)

	logger.Println("Connecting to mongodb..")
	mongo, err := db.New()
	if err != nil {
		logger.Fatalf("Error connecting to mongo: %s\n", err.Error())
		return
	}
	logger.Println("Successfully connected to mongodb")

	tracker := new(sync.WaitGroup)
	router := records.Router(mongo, tracker, logger)

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
	signaler := make(chan os.Signal)
	signal.Notify(signaler, os.Interrupt)
	signal.Notify(signaler, os.Kill)

	<-signaler

	logger.Println("Shutdown signal received, waiting for go routines to finish...")

	tracker.Wait()

	logger.Println("Exiting...")
}
