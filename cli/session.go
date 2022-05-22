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

	connectionScanner := bufio.NewScanner(connection)
	authenticated := authenticateSession(connection, connectionScanner, username, password)
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
			fmt.Println("printed help")
		case "send":
			file, metaData, err := getFileFromMemory(input[1:])
			if err != nil {
				fmt.Printf("Error reading file => %s\n", err.Error())
				break
			}
			sendFileToServer(file, metaData, connection)
		case "get":
		case "quit":
			fmt.Println("Exiting...")
			return
		}
	}
}
