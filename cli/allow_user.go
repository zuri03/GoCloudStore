package cli

import "fmt"

func addAllowedUserCommand(owner string, input []string,
	serverClient RecordServerClient) {
	key := input[0]
	allowedUser := input[1]
	if err := serverClient.AddAllowedUser(owner, key, allowedUser); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	fmt.Println("Successfully added user to allow list")
}
