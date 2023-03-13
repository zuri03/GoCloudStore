package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/zuri03/GoCloudStore/cli"
)

//Read command line arguments
//Attempt to create sesison struct from arguments
//run loop to read input and execute commands
func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Not enough arguments \n Exiting...")
		return
	}

	//first arg should be 'cli' to start a cli session or a command
	firstArg := args[0]
	metadataClient := cli.MetaDataClient{Client: http.Client{Timeout: time.Duration(5) * time.Second}}

	session := cli.ParseArgsIntoStruct(args[2:])

	if firstArg == "cli" {
		cli.HandleSession(&metadataClient)
		return
	} else {
		cli.HandleOneTime(&metadataClient, args)
		return
	}
}
