package cli

import (
	"net"

	c "github.com/zuri03/GoCloudStore/constants"
)

func deleteFile(username string, password string, input []string, metaClient *MetadataServerClient) error {
	key := input[0]
	err := metaClient.deleteFileRecord(username, password, key)
	if err != nil {
		return err
	}

	connection, err := net.Dial("tcp", ":8000")
	defer connection.Close()
	connection.Write([]byte(c.DELETE_PROTOCOL))
	return nil
}
