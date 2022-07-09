package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/zuri03/GoCloudStore/cli"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Not enough arguments")
		return
	}

	firstArg := args[0]
	metadataClient := cli.MetaDataClient{Client: http.Client{Timeout: time.Duration(5) * time.Second}}
	if firstArg == "cli" {
		cli.HandleSession(&metadataClient)
		return
	} else {
		cli.HandleOneTime(&metadataClient, args)
		return
	}
}
