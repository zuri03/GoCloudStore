package cli

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func ExecuteSingleCommand(input []string) {

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
		return
	}
	username := input[1]
	password := input[2]

	metadataClient := MetadataServerClient{Client: http.Client{Timeout: time.Duration(5) * time.Second}}

	switch command {
	case "allow":
		addAllowedUserCommand(username, password, input[3:], &metadataClient)
	case "remove":
		removeUserAccessCommand(username, password, input[3:], &metadataClient)
	case "send":
		sendFileCommand(username, password, input[3:], &metadataClient)
	case "get":
		getFileCommand(username, password, input[3:], &metadataClient)
	case "delete":
		deleteFile(username, password, input[3:], &metadataClient)
	case "help":
		printHelpMessage()
	}
}
