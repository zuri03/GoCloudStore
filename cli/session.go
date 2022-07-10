package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func cleanUserInput(r rune) bool {
	if unicode.IsGraphic(r) {
		return false
	}
	return true
}

func HandleOneTime(client *MetaDataClient, input []string) {
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

	if strings.ToLower(command) == "create" {
		if err := client.createUser(username, password); err != nil {
			fmt.Printf("Error creating user: %s\n", err.Error())
			return
		}
	}

	if _, err := executeCommand(client, command, username, password, input[3:]); err != nil {
		fmt.Print(err.Error())
	}
}

func HandleSession(client *MetaDataClient) {

	commandLineReader := bufio.NewReader(os.Stdin)

	fmt.Printf("username:")
	username, err := commandLineReader.ReadString('\n')
	if err != nil {
		fmt.Printf("Internal Cli error: %s\n Exiting...", err.Error())
		return
	}
	username = strings.TrimFunc(username, cleanUserInput)
	fmt.Printf("password:")

	password, err := commandLineReader.ReadString('\n')
	if err != nil {
		fmt.Printf("Internal Cli error: %s\n Exiting...", err.Error())
		return
	}
	password = strings.TrimFunc(password, cleanUserInput)

	if exists, err := client.authenticate(username, password); !exists {
		fmt.Printf("User %s does not exist. Would you like to create a new user?\n", username)
		fmt.Printf("Yes(Y) or NO(N):")

		response, err := commandLineReader.ReadString('\n')
		if err != nil {
			fmt.Printf("Internal Cli error: %s\n Exiting...", err.Error())
			return
		}
		response = strings.TrimFunc(response, cleanUserInput)

		if strings.ToLower(response) == "y" {
			if err := client.createUser(username, password); err != nil {
				fmt.Printf("Error creating user: %s\n", err.Error())
				return
			}
		}
	} else if err != nil {
		fmt.Printf("Error authorizing user: %s\n", err.Error())
		return
	}

	for {
		fmt.Printf(">")
		str, err := commandLineReader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input => %s\n", err.Error())
			return
		}

		//Trim the two invisible characters at the end
		str = str[:len(str)-2]
		input := strings.Split(str, " ")
		if quit, err := executeCommand(client, input[0], username, password, input[1:]); err != nil || quit {
			fmt.Println(err.Error())
			break
		}
	}

	fmt.Printf("Closing connection")
}

func executeCommand(metadataClient *MetaDataClient, command string, username string, password string, input []string) (bool, error) {

	switch strings.ToLower(command) {
	case "help":
		printHelpMessage()
	case "allow":
		addAllowedUserCommand(username, password, input, metadataClient)
	case "remove":
		removeUserAccessCommand(username, password, input, metadataClient)
	case "send":
		sendFileCommand(username, password, input, metadataClient)
	case "get":
		getFileCommand(username, password, input, metadataClient)
	case "delete":
		deleteFile(username, password, input, metadataClient)
	case "quit":
		fmt.Println("Exiting...")
		return true, nil
	}

	return false, nil
}
