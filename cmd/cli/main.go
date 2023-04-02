package main

import (
	"fmt"
	"os"
	"unicode"

	"github.com/zuri03/GoCloudStore/cli"
	"github.com/zuri03/GoCloudStore/clients"
)

func cleanUserInput(r rune) bool {
	if unicode.IsGraphic(r) && !unicode.IsSpace(r) {
		return false
	}
	return true
}

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

	recordServerClient := clients.NewRecordServerClient()
	fileServerClient := clients.NewFileServerClient()

	session := cli.ParseArgsIntoSession(args[2:])

	if firstArg == "cli" {
		cli.HandleSession(fileServerClient, recordServerClient, session)
	} else {
		cli.HandleOneTime(fileServerClient, recordServerClient, args, session)
	}
}
