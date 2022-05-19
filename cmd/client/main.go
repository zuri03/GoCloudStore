package main

import (
	"fmt"
	"os"

	"github.com/zuri03/GoCloudStore/cli"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Not enough arguments")
		return
	}

	firstArg := args[0]

	if firstArg == "cli" {
		cli.HandleCliSession()
		return
	} else {
		//HandleFileTransfer();
		return
	}

	/*
		if fileExtension != ".txt" && fileExtension != ".rtf" && fileExtension != ".pdf" {
			fmt.Printf("Not a text file")
			return
		}

		file, err := os.Open(fileName)
		if err != nil {
			fmt.Printf("File %s does not exist", fileName)
			return
		}

		buffer := make([]byte, 5)

		connection, err := net.Dial("tcp", ":8080")
		if err != nil {
			fmt.Printf("Error dialing tcp => %s\n", err.Error())
			return
		}
		defer connection.Close()

		fmt.Printf("Client has connected \n")

		for {
			numOfBytes, err := file.Read(buffer)

			if err != nil {
				if err.Error() == "EOF" {
					fmt.Println("End of file found")
					connection.Write([]byte("EOF"))
					return
				} else {
					fmt.Printf("Error occured while reading file => %s\n", err.Error())
					break
				}
			}

			if numOfBytes == 0 {
				fmt.Println("Finished reading file")
				break
			}

			connection.Write(buffer)

			fmt.Printf("=> %s \n", string(buffer[:numOfBytes]))
		}
	*/
}
