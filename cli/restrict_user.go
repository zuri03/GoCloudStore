package cli

import "fmt"

func removeUserAccessCommand(owner string, input []string,
	metaClient *MetaDataClient) error {
	key := input[0]
	removedUser := input[1]
	if err := metaClient.removeAllowedUser(owner, key, removedUser); err != nil {
		return err
	}
	fmt.Printf("Successfully removed %s from allow list\n", removedUser)
	return nil
}
