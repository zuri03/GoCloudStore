package server

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

//TODO: Turn all of the buffer sizes into contants
//TODO: Determine proper buffer size
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

type FileMetaData struct {
	Size int64
	Name string
}

func acceptFileMetaData(connectionScanner *bufio.Scanner, connection net.Conn) (*FileMetaData, error) {
	gob.Register(new(FileMetaData))
	meta := &FileMetaData{}
	decoder := gob.NewDecoder(connection)
	fmt.Println("WAITING FOR GOB")
	err := decoder.Decode(meta)
	if err != nil {
		fmt.Printf("Error decoding gob => %s\n", err.Error())
		return nil, err
	}
	return meta, nil
}

func createFile(metaData *FileMetaData) (*os.File, error) {
	file, err := os.OpenFile(metaData.Name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return nil, err
	}
	fmt.Println("CREATED FILE")
	return file, nil
}

func acceptFileData(metaData *FileMetaData, file *os.File, connection net.Conn) error {
	if metaData.Size >= 1024 {
		//Read the data from the connection in a loop
	} else {
		dataBuffer := make([]byte, 1024)
		connection.Read(dataBuffer)
		result := trimBuffer(dataBuffer)
		file.Write(result)
	}
	return nil
}

func storeFileFromClient(connectionScanner *bufio.Scanner, connection net.Conn) error {
	fmt.Println("ACCEPTING META DATA")
	metaData, err := acceptFileMetaData(connectionScanner, connection)
	if err != nil {
		return err
	}

	fmt.Println("CREATE FILE")
	file, err := createFile(metaData)
	if err != nil {
		return err
	}

	fmt.Println("ACCEPT FILE DATA")
	err = acceptFileData(metaData, file, connection)
	if err != nil {
		return err
	}

	fmt.Println("SUCCESS")
	return nil
}

func sendFileToClient(fileName string, connection net.Conn) error {
	metaData, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	fmt.Printf("size => %d\n", metaData.Size())
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	if metaData.Size() >= 1024 {
		buffer := make([]byte, metaData.Size())
		file.Read(buffer)
		connection.Write(buffer)
		fmt.Println("SENT DATA TO CLIENT")
	}

	return nil
}

func HandleConnection(connection net.Conn) {
	fmt.Println("Handling connection")
	defer connection.Close()

	//Remove connection scanner
	connectionScanner := bufio.NewScanner(connection)
	authenticated := authenticateConnection(connectionScanner, connection)

	fmt.Println("exited auth")
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
			err := storeFileFromClient(connectionScanner, connection)
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

func trimBuffer(buffer []byte) []byte {
	fmt.Printf("TRIMMING => %s\n", string(buffer))
	return bytes.TrimFunc(buffer, func(r rune) bool {
		if r == 0 {
			return true
		}
		return false
	})
}
