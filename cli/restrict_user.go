package cli

func removeUserAccessCommand(username string, password string, input []string,
	metaClient *MetaDataClient) error {
	key := input[0]
	removedUser := input[1]
	if err := metaClient.removeAllowedUser(username, password, key, removedUser); err != nil {
		return err
	}
	return nil
}
