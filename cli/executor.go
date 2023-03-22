package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func CleanUserInput(r rune) bool {
	if unicode.IsGraphic(r) && !unicode.IsSpace(r) {
		return false
	}
	return true
}

func authenticateSession(metadataClient *MetaDataClient, session *Session, commandLineReader *bufio.Reader) error {
	id, exists, err := metadataClient.authenticate(session.Username, session.Password)
	if err != nil {
		return err
	}

	if !exists {
		fmt.Printf("User %s does not exist. Would you like to create a new user?\n", session.Username)
		fmt.Printf("Yes(Y) or NO(N):")

		response, err := commandLineReader.ReadString('\n')
		if err != nil {
			return err
		}
		response = strings.TrimFunc(response, CleanUserInput)

		if strings.ToLower(response) == "y" {
			user, err := metadataClient.createUser(session.Username, session.Password)
			if err != nil {
				fmt.Printf("Error creating user: %s\n", err.Error())
				return err
			}
			id = user.Id
		} else {
			return errors.New("User credentials were incorrect\n")
		}
	}

	session.Id = id

	return nil
}

func HandleOneTime(client *MetaDataClient, input []string, session *Session) {
	if len(input) == 0 {
		fmt.Println("Incorrect number of arguments")
		printHelpMessage()
		return
	}

	command := strings.ToLower(input[0])
	if command == "help" {
		printHelpMessage()
		return
	}

	if err := authenticateSession(client, session, bufio.NewReader(os.Stdin)); err != nil {
		fmt.Print(err.Error())
		return
	}

	if _, err := executeCommand(client, command, session.Id, input[3:]); err != nil {
		fmt.Print(err.Error())
	}
}

func HandleSession(client *MetaDataClient, session *Session) {

	commandLineReader := bufio.NewReader(os.Stdin)

	if err := authenticateSession(client, session, commandLineReader); err != nil {
		fmt.Print(err.Error())
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
		if quit, err := executeCommand(client, input[0], session.Id, input[1:]); err != nil || quit {
			if err != nil {
				fmt.Println(err.Error())
			}
			break
		}
	}

	fmt.Printf("Exiting...")
}

func executeCommand(metadataClient *MetaDataClient, command string, owner string, input []string) (bool, error) {
	fmt.Printf("Command => %s\n", command)
	switch strings.ToLower(command) {
	case "help":
		printHelpMessage()
	case "allow":
		addAllowedUserCommand(owner, input, metadataClient)
	case "remove":
		removeUserAccessCommand(owner, input, metadataClient)
	case "send":
		sendFileCommand(owner, input, metadataClient)
	case "get":
		getFileCommand(owner, input, metadataClient)
	case "delete":
		deleteFile(owner, input, metadataClient)
	case "quit":
		fmt.Println("Exiting...")
		return true, nil
	default:
		fmt.Println("Unrecognized Command")
	}

	return false, nil
}
