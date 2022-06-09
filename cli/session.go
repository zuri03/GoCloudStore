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
			err := addAllowedUserCommand(username, password, input[1:], &metadataClient)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("Successfully added user to permissions list")
		case "remove":
			err := removeUserAccessCommand(username, password, input[1:], &metadataClient)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("Successfully removed user from permissions list")
		case "send":
			err := sendFileCommand(username, password, input[1:], &metadataClient)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("Successfully sent file to server")
		case "get":
			err := getFileCommand(username, password, input[1:], &metadataClient)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("Successuflly retreived file from server")
		case "delete":
			err := deleteFile(username, password, input[1:], &metadataClient)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("Successuflly deleted file from server")
		case "quit":
			fmt.Println("Exiting...")
			return
		}
	}
}
