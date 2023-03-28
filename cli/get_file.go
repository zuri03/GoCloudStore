package cli

import (
	"fmt"
)

//Opens and closes connections and calls the file retreival functions
//in the correct order
//Handles any errors that may occur by notifying the user
func getFileCommand(owner string, input []string, serverClient ServerClient) {
	key := input[0]
	record, err := serverClient.GetFileRecord(owner, key)
	if err != nil {
		fmt.Printf("Error getting file record: %s\n", err.Error())
		return
	}

	if _, err := serverClient.GetFile(owner, record); err != nil {
		fmt.Println(err.Error())
	}
}
