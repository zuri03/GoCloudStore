package cli

import (
	"bufio"
	//"encoding/json"
	"encoding/gob"
	"fmt"
	"io/fs"
	"net"
	"os"
)

func authenticateSession(connection net.Conn, scanner *bufio.Scanner, username string, password string) bool {
	fmt.Println("Authenticating")
	connection.Write([]byte(username))
	connection.Write([]byte(password))
	buf := make([]byte, 2)
	connection.Read(buf)
	if string(buf) != "OK" {
		return false
	}
	fmt.Println("SUCCESSFUL AUTHENTICATION")
	return true
}

type FileMetaData struct {
	Size int64
	Name string
}

//Place types in a shared directory in the future
func sendMetaDataToServer(meta fs.FileInfo, connection net.Conn) error {

	fmt.Println("Generated gob")
	/*
		jsonBytes, err := json.Marshal(meta)
		if err != nil {
			return err
		}
		connection.Write(jsonBytes)
		connection.Write([]byte("\n"))

		fmt.Println("Sent json")
	*/

	gob.Register(new(FileMetaData))
	metaData := FileMetaData{
		Size: meta.Size(),
		Name: meta.Name(),
	}
	encoder := gob.NewEncoder(connection)
	fmt.Println("Connected gob to buffer")
	err := encoder.Encode(metaData)
	if err != nil {
		fmt.Printf("Error in gob => %s\n", err.Error())
		return err
	}
	fmt.Println("Encoded meta data")
	fmt.Println("Encoded meta data")
	return nil
}

func sendDataToServer() {

}

func sendFileToServer(file *os.File, meta fs.FileInfo, connection net.Conn) error {

	fmt.Println("Sending meta data to server")
	err := sendMetaDataToServer(meta, connection)
	if err != nil {
		return err
	}

	return nil

	/*

		//Arbitrary buffer size
		//Determine best buffer size later
		buffer := make([]byte, 500)

		for {
			numOfBytes, err := file.Read(buffer)

			if err != nil {
				if err.Error() == "EOF" {
					fmt.Println("End of file found")
					connection.Write([]byte("EOF"))
					return nil
				} else {
					return fmt.Errorf("Error occured while reading file => %s\n", err.Error())
				}
			}

			fmt.Printf("Number of bytes => %d\n", numOfBytes)
			if numOfBytes == 0 {
				fmt.Println("Finished reading file")
				break
			}

			connection.Write(buffer)
		}
		return nil
	*/
}
