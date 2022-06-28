package main

import (
	"fmt"
	"os"

	"github.com/zuri03/GoCloudStore/cli"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Not enough arguments")
		return
	}

	firstArg := args[0]

	if firstArg == "cli" {
		cli.HandleCliSession()
		return
	} else {
		cli.ExecuteSingleCommand(args)
		//HandleFileTransfer();
		return
	}
}
