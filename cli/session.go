package cli

import (
	"bufio"
	"fmt"
	"net"
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
		case "send":
			file, meta, err := getFileFromMemory(input[1:])
			if err != nil {
				fmt.Println(err.Error())
			}
			if meta == nil || file == nil {
				fmt.Println("Error file not found")
				continue
			}
			err = metadataClient.createFileRecord(username, password, meta.Name(), meta.Name(), meta.Size()) //For now just leave the key as the file name
			if err != nil {
				fmt.Println(err.Error())
				continue
			} else {
				fmt.Println("NO ERROR")
			}
		case "get":
			record, err := metadataClient.getFileRecord(username, password, input[1])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Printf("record => %s\n", record.MetaData.Name)
		case "delete":
			deleteFile(username, password, input[1:], &metadataClient)
		case "quit":
			fmt.Println("Exiting...")
			return
		}
	}
}
