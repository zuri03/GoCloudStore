package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//This struct will hold a set of values that remain constant throughout a single cli session
//The struct is recreated everytime the program is run
type Session struct {
	Username string
	Password string
	Id       string
}

//Accepts an array of command line arguments and attempts to create a session object from it
//This function will assume that the arguments are in this format [ <username>, <password> ]
//If it is
func ParseArgsIntoSession(commandLineArgs []string) Session {
	newSession := Session{}

	if len(commandLineArgs) < 2 {
		fillSessionStruct(&newSession)
		return newSession
	}

	username := commandLineArgs[0]
	password := commandLineArgs[1]

	newSession.Username = strings.TrimFunc(username, CleanUserInput)
	newSession.Password = strings.TrimFunc(password, CleanUserInput)

	return newSession
}

//Ask the user for any missing fields in the session struct
func fillSessionStruct(session *Session) error {

	reader := bufio.NewReader(os.Stdin)

	if session.Username == "" {
		fmt.Printf("Please enter your username: ")

		response, err := reader.ReadString('\n')

		if err != nil {
			return err
		}

		session.Username = response
	}

	if session.Password == "" {
		fmt.Printf("Please enter your password: ")

		response, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		session.Password = response
	}

	return nil
}
