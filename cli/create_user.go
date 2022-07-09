package cli

import (
	"fmt"
)

func createUesr(username string, password string, metaClient *MetaDataClient) (string, string, error) {

	err := metaClient.create(username, password)
	if err != nil {
		fmt.Printf("Error retreiving meta data from server: %s\n", err.Error())
		return "", "", nil
	}

	fmt.Println("Successfully creatd user")
	return username, password, nil
}
