package cli

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func ExecuteSingleCommand(input []string) {

	fmt.Println("EXECUTING ONE TIME COMMAND")


	if len(input) == 0 {
		fmt.Println("Error: No arguments found ")
		return
	}

	command := strings.ToLower(input[0])
	if command == "help" {
		printHelpMessage()
		return
	}

	if len(input) < 3 {
		fmt.Println("Incorrect number of arguments. Correct format \n \t" +
			"GoCloudStore [command] [username] [password] [command arguments]")
		fmt.Println("")
		return
	}
	username := input[1]
	password := input[2]

	metadataClient := MetadataServerClient{Client: http.Client{Timeout: time.Duration(5) * time.Second}}

	switch command {
	case "allow":
		err := addAllowedUserCommand(username, password, input[3:], &metadataClient)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Successfully added user to permissions list")
	case "remove":
		err := removeUserAccessCommand(username, password, input[3:], &metadataClient)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Successfully removed user from permissions list")
	case "send":
		err := sendFileCommand(username, password, input[3:], &metadataClient)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Successfully sent file to server")
	case "get":
		err := getFileCommand(username, password, input[3:], &metadataClient)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Successuflly retreived file from server")
	case "delete":
		err := deleteFile(username, password, input[3:], &metadataClient)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Successuflly deleted file from server")
	case "quit":
		fmt.Println("Exiting...")
		return
	}
}
