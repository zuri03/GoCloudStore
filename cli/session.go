package cli

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
	authenticateSession(connection, connectionScanner, username, password)
	fmt.Printf("Closing connection")
}
