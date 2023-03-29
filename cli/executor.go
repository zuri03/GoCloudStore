package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"unicode"

	"github.com/zuri03/GoCloudStore/common"
)

type RecordServerClient interface {
	AuthenticateUser(username string, password string) (string, bool, error)
	CreateUser(username string, password string) (string, error)
	GetFileRecord(owner string, key string) (*common.Record, error)
	DeleteFileRecord(owner string, key string) error
	CreateFileRecord(owner, key, fileName string, fileSize int64) error
	AddAllowedUser(owner, key, allowedUser string) error
	RemoveAllowedUser(owner, key string, removedUser string) error
}

type FileServerClient interface {
	SendFile(owner string, file *os.File, fileInfo fs.FileInfo) error
	GetFile(owner string, record *common.Record) (string, error)
	DeleteFile(owner string, record *common.Record) error
}

func CleanUserInput(r rune) bool {
	if unicode.IsGraphic(r) && !unicode.IsSpace(r) {
		return false
	}
	return true
}

func authenticateSession(serverClient RecordServerClient, session Session, commandLineReader *bufio.Reader) error {
	id, exists, err := serverClient.AuthenticateUser(session.Username, session.Password)
	if err != nil {
		return err
	}

	if !exists {
		fmt.Printf("User %s does not exist. Would you like to create a new user?\n", session.Username)
		fmt.Printf("Yes(Y) or NO(N):")

		response, err := commandLineReader.ReadString('\n')
		if err != nil {
			return err
		}
		response = strings.TrimFunc(response, CleanUserInput)

		if strings.ToLower(response) == "y" {
			id, err = serverClient.CreateUser(session.Username, session.Password)
			if err != nil {
				fmt.Printf("Error creating user: %s\n", err.Error())
				return err
			}
		} else {
			return errors.New("User credentials were incorrect\n")
		}
	}

	session.Id = id

	return nil
}

func HandleOneTime(fileServerClient FileServerClient, serverClient RecordServerClient, input []string, session Session) {
	if len(input) == 0 {
		fmt.Println("Incorrect number of arguments")
		printHelpMessage()
		return
	}

	command := strings.ToLower(input[0])
	if command == "help" {
		printHelpMessage()
		return
	}

	if err := authenticateSession(serverClient, session, bufio.NewReader(os.Stdin)); err != nil {
		fmt.Print(err.Error())
		return
	}

	if _, err := executeCommand(fileServerClient, serverClient, command, session.Id, input[3:]); err != nil {
		fmt.Print(err.Error())
	}
}

func HandleSession(fileServerClient FileServerClient, serverClient RecordServerClient, session Session) {

	commandLineReader := bufio.NewReader(os.Stdin)

	if err := authenticateSession(serverClient, session, commandLineReader); err != nil {
		fmt.Print(err.Error())
		return
	}

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
		if quit, err := executeCommand(fileServerClient, serverClient, input[0], session.Id, input[1:]); err != nil || quit {
			if err != nil {
				fmt.Println(err.Error())
			}
			break
		}
	}

	fmt.Printf("Exiting...")
}

func executeCommand(fileServerClient FileServerClient, serverClient RecordServerClient, command string, owner string, input []string) (bool, error) {
	fmt.Printf("Command => %s\n", command)
	switch strings.ToLower(command) {
	case "help":
		printHelpMessage()
	case "allow":
		addAllowedUserCommand(owner, input, serverClient)
	case "remove":
		removeUserAccessCommand(owner, input, serverClient)
	case "send":
		sendFileCommand(owner, input, serverClient, fileServerClient)
	case "get":
		getFileCommand(owner, input, serverClient, fileServerClient)
	case "delete":
		deleteFile(owner, input, serverClient, fileServerClient)
	case "quit":
		fmt.Println("Exiting...")
		return true, nil
	default:
		fmt.Println("Unrecognized Command")
	}

	return false, nil
}
