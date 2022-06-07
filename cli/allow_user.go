package cli

func addAllowedUserCommand(username string, password string, input []string,
	metaClient *MetadataServerClient) error {
	key := input[0]
	allowedUser := input[1]
	if err := metaClient.addAllowedUser(username, password, key, allowedUser); err != nil {
		return err
	}
	return nil
}
