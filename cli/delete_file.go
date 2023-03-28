package cli

import (
	"fmt"
)

func deleteFile(owner string, input []string, serverClient ServerClient) {
	key := input[0]

	record, err := serverClient.GetFileRecord(owner, key)
	if err != nil {
		fmt.Printf("Error retreiving meta data from server: %s\n", err.Error())
		return
	}

	err = serverClient.DeleteFileRecord(owner, key)
	if err != nil {
		fmt.Printf("Error deleting meta data from server: %s\n", err.Error())
		return
	}

	if err := serverClient.DeleteFile(owner, record); err != nil {
		fmt.Println("Error deleting file off of file server: %s\n", err.Error())
	}
}
