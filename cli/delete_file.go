package cli

import (
	"fmt"
	"net"
)

func deleteFile(username string, password string, input []string, metaClient *MetadataServerClient) error {
	key := input[0]
	fmt.Printf("key => %s\n", key)
	err := metaClient.deleteFileRecord(username, password, key)
	if err != nil {
		return err
	}

	connection, err := net.Dial("tcp", ":8000")
	defer connection.Close()
	fmt.Println("ABOUT TO GET")
	connection.Write([]byte(DELETE_PROTOCOL))
	return nil
}
