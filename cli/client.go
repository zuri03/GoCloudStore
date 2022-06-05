package cli

import (
	//"encoding/json"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/zuri03/GoCloudStore/records"
)

func authenticateSession(connection net.Conn, username string, password string) bool {
	fmt.Println("Authenticating")
	connection.Write([]byte(username))
	connection.Write([]byte(password))
	buf := make([]byte, 2)
	connection.Read(buf)
	if string(buf) != "OK" {
		return false
	}
	fmt.Println("SUCCESSFUL AUTHENTICATION")
	return true
}

type FileMetaData struct {
	Size int64
	Name string
}

//Place types in a shared directory in the future
func sendMetaDataToServer(meta fs.FileInfo, connection net.Conn) error {
	fmt.Println("Generated gob")
	gob.Register(new(FileMetaData))
	metaData := FileMetaData{
		Size: meta.Size(),
		Name: meta.Name(),
	}
	encoder := gob.NewEncoder(connection)
	fmt.Println("Connected gob to buffer")
	err := encoder.Encode(metaData)
	if err != nil {
		return err
	}
	fmt.Println("Encoded meta data")
	return nil
}

func sendFileDataToServer(file *os.File, meta fs.FileInfo, connection net.Conn) error {
	fmt.Printf("FILE SIZE => %d\n", meta.Size())
	if meta.Size() >= 1024 {
		buffer := make([]byte, 1024)

		for {
			numOfBytes, err := file.Read(buffer)

			if err != nil {
				if err.Error() == "EOF" {
					fmt.Println("End of file found")
					connection.Write([]byte("EOF"))
					return nil
				} else {
					return fmt.Errorf("Error occured while reading file => %s\n", err.Error())
				}
			}

			fmt.Printf("Number of bytes => %d\n", numOfBytes)
			if numOfBytes == 0 {
				fmt.Println("Finished reading file")
				break
			}

			connection.Write(buffer)
		}
	} else {
		fmt.Println("SENDING FILE IN ONE CHUNK")
		dataBuffer := make([]byte, meta.Size())
		file.Read(dataBuffer)
		connection.Write(dataBuffer)
		fmt.Println("SENT FILE DATA")
	}
	return nil
}

func sendFileToServer(file *os.File, meta fs.FileInfo, connection net.Conn) error {

	fmt.Println("Sending meta data to server")
	err := sendMetaDataToServer(meta, connection)
	if err != nil {
		return err
	}
	fmt.Println("SENDING FILE DATA TO SERVER")
	sendFileDataToServer(file, meta, connection)
	return nil
}

type MetadataServerClient struct {
	Client http.Client
}

//Meta data server functions
func (c *MetadataServerClient) getFileRecord(username string, password string, key string) (*records.Record, error) {

	url := fmt.Sprintf("http://localhost:8080/record?username=%s&password=%s&key=%s", username, password, key)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(request)
	if err != nil {
		return nil, err
	}

	var record records.Record
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error has occured while getting record: %s\n", string(body))
	}

	json.Unmarshal(body, &record)
	return &record, nil
}

func (c *MetadataServerClient) deleteFileRecord(username string, password string, key string) error {

	url := fmt.Sprintf("http://localhost:8080/record?username=%s&password=%s&key=%s", username, password, key)
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(request)
	if err != nil {
		return err
	}

	var record records.Record
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Error has occured while deleting record: %s\n", string(body))
	}

	json.Unmarshal(body, &record)
	return nil
}

func (c *MetadataServerClient) createFileRecord(username string, password string, key string, fileName string, fileSize int64) error {

	record := records.CreateReqest{
		Username: username,
		Password: password,
		Key:      key,
		FileName: fileName,
		Size:     fileSize,
	}

	recordBytes, _ := json.Marshal(record)
	request, err := http.NewRequest("POST", "http://localhost:8080/record", bytes.NewBuffer(recordBytes))
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		errorMessage, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error has occured while creating record: %s\n", string(errorMessage))
	}

	return nil
}
