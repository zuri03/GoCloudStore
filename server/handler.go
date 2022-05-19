package server

import (
	"bufio"
	"fmt"
	"net"
	//"os"
)

func authenticateConnection(connectionScanner *bufio.Scanner, connection net.Conn) bool {
	fmt.Println("AUTHENTICATING CLIENT")
	connectionScanner.Scan()
	fmt.Println("GOT TOKEN")

	username := connectionScanner.Text()
	fmt.Printf("Username => %s\n", username)
	connectionScanner.Scan()
	//connection.Write([]byte("OK"))

	password := connectionScanner.Text()
	fmt.Printf("Password => %s\n", password)

	//From here pass username and password to some authenticaiton service
	fmt.Println("SUCCESS ENDING SESSION")
	connection.Write([]byte("OK"))
	return true
}
func HandleConnection(connection net.Conn) {
	fmt.Println("Handling connection")
	defer connection.Close()

	connection.Write([]byte("OK"))

	fmt.Println("SENT OK")

	connectionScanner := bufio.NewScanner(connection)

	authenticated := authenticateConnection(connectionScanner, connection)

	if !authenticated {
		connection.Write([]byte("User Not Found"))
		return
	}

	/*
		file, err := os.OpenFile("examle.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			fmt.Printf("Error occured while opening file => %s\n", err.Error())
			return
		}

		for connectionScanner.Scan() {
			text := connectionScanner.Text()

			fmt.Printf("Line => %s\n", text)

			if text == "EOF" {
				fmt.Println("End of file received")
				return
			}

			file.WriteString(text)
		}
	*/
}
