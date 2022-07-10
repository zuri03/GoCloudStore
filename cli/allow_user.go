package cli

import "fmt"

func addAllowedUserCommand(username string, password string, input []string,
	metaClient *MetaDataClient) error {
	key := input[0]
	allowedUser := input[1]
	if err := metaClient.addAllowedUser(username, password, key, allowedUser); err != nil {
		return err
	}
	fmt.Println("Successfully added user to allow list")
	return nil
}
