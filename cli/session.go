package cli

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
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
	runSessionLoop(commandLineReader, connection, username, password)
	fmt.Printf("Closing connection")
}

func runSessionLoop(commandLineReader *bufio.Reader, connection net.Conn, username string, password string) {
	metadateClient := MetadataServerClient{Client: http.Client{Timeout: time.Duration(5) * time.Second}}
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
		case "get":
			record, err := metadateClient.getFileRecord(username, password, input[1])
			if err != nil {
				fmt.Printf("Error => %s\n", err.Error())
			}
			fmt.Printf("record => %s\n", record.MetaData.Name)
		case "quit":
			fmt.Println("Exiting...")
			return
		}
	}
}
