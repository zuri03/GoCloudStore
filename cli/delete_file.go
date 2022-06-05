package cli

import "fmt"

func deleteFile(username string, password string, input []string, metaClient *MetadataServerClient) error {
	key := input[0]
	fmt.Printf("key => %s\n", key)
	err := metaClient.deleteFileRecord(username, password, key)
	if err != nil {
		return err
	}
	return nil
}
