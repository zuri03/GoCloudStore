package cli

import (
	"bufio"
	"fmt"
	"net"
)

func getServerResponse(scanner *bufio.Scanner) (string, error) {
	fmt.Printf("Waiting for message")
	scanner.Scan()
	message := scanner.Text()
	fmt.Printf("Connection successful message => %s\n", message)

	if message != "OK" {
		return "", fmt.Errorf("Error expected OK got => %s\n", message)
	}

	return message, nil
}

func authenticateSession(connection net.Conn, scanner *bufio.Scanner, username string, password string) bool {

	fmt.Println("Authenticating")
	//Make sure connection was successful
	/*
		_, err := getServerResponse(scanner)
		if err != nil {
			fmt.Printf("Error connecting: %s\n", err.Error())
			return false
		}
	*/

	connection.Write([]byte(username))

	/*
		_, err = getServerResponse(scanner)
		if err != nil {
			fmt.Printf("Error sending username: %s\n", err.Error())
			return false
		}
	*/

	connection.Write([]byte(password))

	/*
		_, err = getServerResponse(scanner)
		if err != nil {
			fmt.Printf("Error connecting: %s\n", err.Error())
			return false
		}
	*/

	scanner.Scan()
	response := scanner.Text()
	fmt.Printf("Response => %s\n", response)

	fmt.Println("SUCCESSFUL AUTHENTICATION")

	return false
}
