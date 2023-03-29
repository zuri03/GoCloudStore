package cli

import (
	"fmt"
)

func deleteFile(owner string, input []string, recordServerClient RecordServerClient, fileServerClient FileServerClient) {
	key := input[0]

	record, err := recordServerClient.GetFileRecord(owner, key)
	if err != nil {
		fmt.Printf("Error retreiving meta data from server: %s\n", err.Error())
		return
	}

	err = recordServerClient.DeleteFileRecord(owner, key)
	if err != nil {
		fmt.Printf("Error deleting meta data from server: %s\n", err.Error())
		return
	}

	if err := fileServerClient.DeleteFile(owner, record); err != nil {
		fmt.Printf("Error deleting file off of file server: %s\n", err.Error())
	}
}
