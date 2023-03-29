package cli

import "fmt"

func removeUserAccessCommand(owner string, input []string,
	serverClient RecordServerClient) error {
	key := input[0]
	removedUser := input[1]
	if err := serverClient.RemoveAllowedUser(owner, key, removedUser); err != nil {
		return err
	}
	fmt.Printf("Successfully removed %s from allow list\n", removedUser)
	return nil
}
