package server

import (
	"bufio"
	"bytes"
	"fmt"
	"net"

	"github.com/zuri03/GoCloudStore/records"
)

//TODO: Turn all of the buffer sizes into contants
//TODO: Determine proper buffer size
//TODO: Move storage functionality to the storage server
type Handler struct {
	Keeper records.RecordKeeper
}

func (h *Handler) HandleConnection(connection net.Conn) {
	defer connection.Close()

	//Remove connection scanner
	connectionScanner := bufio.NewScanner(connection)
	authenticated := authenticateConnection(connectionScanner, connection)

	if !authenticated {
		connection.Write([]byte("UNF")) //user not found
		return
	}

	//This limits the max length of a file name to 96
	//Client message => COMMAND:PARAMETERS
	//Command is three bytes
	readBuffer := make([]byte, 100)
	for {
		fmt.Println("WAITING FOR MESSAGE")
		_, err := connection.Read(readBuffer)
		if err != nil {
			fmt.Println("ERROR CLOSING CONNECTION")
			return
		}

		switch string(readBuffer[0:3]) {
		case "GET":
			fmt.Println("FOUND GET REQEUST")
			fileName := trimBuffer(readBuffer[4:])
			sendFileToClient(string(fileName), connection)
		case "SND":
			fmt.Println("FOUND SEND REQUEST")
			err := storeFileFromClient(connectionScanner, connection, h.Keeper)
			if err != nil {
				fmt.Printf("ERROR => %s\n", err.Error())
				connection.Write([]byte("ERR"))
				return
			}
		case "ERR":
			fmt.Println("CLIENT ERROR OCCURED")
			return
		default:
			fmt.Println("UNRECOGNIZED REQUEST")
		}
	}
}

func checkUserCredentials(username string, password string) {
	//stub
}

func authenticateConnection(connectionScanner *bufio.Scanner, connection net.Conn) bool {
	connectionScanner.Scan()

	username := connectionScanner.Text()
	connectionScanner.Scan()

	password := connectionScanner.Text()

	//From here pass username and password to some authenticaiton service
	checkUserCredentials(username, password)
	connection.Write([]byte("OK"))
	return true
}

func trimBuffer(buffer []byte) []byte {
	return bytes.TrimFunc(buffer, func(r rune) bool {
		if r == 0 {
			return true
		}
		return false
	})
}
