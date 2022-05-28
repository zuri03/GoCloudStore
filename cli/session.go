package cli

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func HandleCliSession() {

	commandLineReader := bufio.NewReader(os.Stdin)

	fmt.Printf("username:")
	username, err := commandLineReader.ReadString('\n')
	if err != nil {
		fmt.Printf("Internal Cli error: %s\n Exiting...", err.Error())
		return
	}

	fmt.Printf("password:")

	password, err := commandLineReader.ReadString('\n')
	if err != nil {
		fmt.Printf("Internal Cli error: %s\n Exiting...", err.Error())
		return
	}

	//Replace with a proper ip address later
	connection, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Printf("Error connecting => %s\n", err.Error())
		return
	}
	defer connection.Close()

	//connectionScanner := bufio.NewScanner(connection)
	authenticated := authenticateSession(connection, username, password)
	if !authenticated {
		fmt.Printf("Failed to authenticate user \n Exiting...")
		return
	}
	runSessionLoop(commandLineReader, connection)
	fmt.Printf("Closing connection")
}

func runSessionLoop(commandLineReader *bufio.Reader, connection net.Conn) {
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
		case "send":
			err := sendFileCommand(input[1:], connection)
			if err != nil {
				fmt.Printf("Error => %s\n", err.Error())
				break
			}
		case "get":
			err := getFileCommand(input[1:], connection)
			if err != nil {
				fmt.Printf("Error retreiving file from server => %s\n", err.Error())
				break
			}
			fmt.Println("FINISHED RETRIEVING FILE FROM SERVER")
		case "quit":
			fmt.Println("Exiting...")
			return
		}
	}
}
