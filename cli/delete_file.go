package cli

import (
	"encoding/gob"
	"fmt"
	"net"

	"github.com/zuri03/GoCloudStore/common"
)

func deleteFile(owner string, input []string, metaClient *MetaDataClient) {
	key := input[0]

	record, err := metaClient.getFileRecord(owner, key)
	if err != nil {
		fmt.Printf("Error retreiving meta data from server: %s\n", err.Error())
		return
	}

	err = metaClient.deleteFileRecord(owner, key)
	if err != nil {
		fmt.Printf("Error deleting meta data from server: %s\n", err.Error())
		return
	}
	meta := common.FileMetaData{
		Owner: owner,
		Name:  record.Name,
		Size:  record.Size,
	}

	connection, err := net.Dial("tcp", ":8000")
	defer connection.Close()
	encoder := gob.NewEncoder(connection)

	if err := sendMetaDataToServer(common.DELETE_FRAME, meta, encoder); err != nil {
		fmt.Printf("Error sending meta data: %s\n", err.Error())
		return
	}

	signal := make([]byte, 3)
	if _, err := connection.Read(signal); err != nil {
		fmt.Printf("Error occured on storage server: %s\n", err.Error())
		return
	}

	if string(signal) != common.SUCCESS_PROTOCOL {
		fmt.Println("Error on server")
		return
	}
	fmt.Println("Successfully deleted file from server")
}
