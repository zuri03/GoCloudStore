package cli

import (
	"fmt"
	"net"
	"os"
)

const (
	FILE_BUFFER_SIZE int = 1024 //TODO: Determine caching buffer size
)

//Opens and closes connections and calls the file retreival functions
//in the correct order
func getFileCommand(username string, password string, input []string, metaClient *MetadataServerClient) error {
	key := input[0]
	record, err := metaClient.getFileRecord(username, password, key)
	if err != nil {
		return err
	}
	fmt.Printf("record => %+v\n", record)
	//connection, err := net.DialTCP("tcp", nil, dataNodeAddress)
	connection, err := net.Dial("tcp", ":8000")
	defer connection.Close()
	if err := getFileDataFromServer(record.MetaData.Name, int(record.MetaData.Size),
		connection); err != nil {
		return nil
	}
	return nil
}

//uses a connection object to retrieve byte data from a storage server and store it in a file
func getFileDataFromServer(fileName string, fileSize int, connection net.Conn) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	connection.Write([]byte(GET_PROTOCOL))
	if fileSize <= FILE_BUFFER_SIZE {
		buffer := make([]byte, fileSize)
		numOfBytes, err := connection.Read(buffer)
		if err != nil {
			fmt.Printf("ERROR NUM OF BYTES => %d\n", numOfBytes)
			return err
		}
		fmt.Printf("final buffer => %s\n", string(buffer))
		if _, err := file.Write(buffer); err != nil {
			fmt.Printf("ERROR WRITING TO FILE +> %s\n", err.Error())
			return err
		}

		return nil
	}

	//Here we would use buffered io to read larger files
	//or maybe a for loop
	return nil
}
