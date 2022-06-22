package cli

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"
)

func cleanUserInput(r rune) bool {
	if unicode.IsGraphic(r) {
		return false
	}
	return true
}

func HandleCliSession() {

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

	//connectionScanner := bufio.NewScanner(connection)
	runSessionLoop(commandLineReader, username, password)
	fmt.Printf("Closing connection")
}

func runSessionLoop(commandLineReader *bufio.Reader, username string, password string) {
	metadataClient := MetadataServerClient{Client: http.Client{Timeout: time.Duration(5) * time.Second}}
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

		switch strings.ToLower(input[0]) {
		case "help":
			printHelpMessage()
		case "allow":
			addAllowedUserCommand(username, password, input[1:], &metadataClient)
		case "remove":
			removeUserAccessCommand(username, password, input[1:], &metadataClient)
		case "send":
			sendFileCommand(username, password, input[1:], &metadataClient)
		case "get":
			getFileCommand(username, password, input[1:], &metadataClient)
		case "delete":
			deleteFile(username, password, input[1:], &metadataClient)
		case "quit":
			fmt.Println("Exiting...")
			return
		}
	}
}
