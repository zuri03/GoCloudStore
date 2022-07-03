package cli

import (
	//"encoding/json"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	c "github.com/zuri03/GoCloudStore/common"
	"github.com/zuri03/GoCloudStore/records"
)

func newEncoderDecorder(connection net.Conn) (*gob.Encoder, *gob.Decoder) {
	gob.Register(new(c.ProtocolFrame))
	gob.Register(new(FileMetaData))

	return gob.NewEncoder(connection), gob.NewDecoder(connection)
}

//This function forces the client to wait for the server to send a message letting
//the client know it is ready to move to the next step of the process
func waitForProceed(decoder *gob.Decoder) error {
	var frame c.ProtocolFrame
	if err := decoder.Decode(&frame); err != nil {
		return err
	}

	if frame.Type == c.ERROR_FRAME {
		return fmt.Errorf("Error on server: %s\n", string(frame.Data))
	}

	if frame.Type != c.PROCEED_FRAME {
		return fmt.Errorf("Unexpected frame got: %d\n", frame.Type)
	}

	fmt.Printf("Got proceed frame => %d\n", frame.Type)
	return nil
}

func sendMetaDataToServer(frameType c.FrameType, meta FileMetaData, encoder *gob.Encoder) error {
	metaBuffer := new(bytes.Buffer)
	if err := gob.NewEncoder(metaBuffer).Encode(meta); err != nil {
		return err
	}
	fmt.Printf("encoding meta => %+v\n", meta)
	fmt.Printf("encoding meta length=> %d\n", metaBuffer.Len())
	frame := c.ProtocolFrame{
		Type:          frameType,
		PayloadLength: int64(metaBuffer.Len()),
		Data:          metaBuffer.Bytes(),
	}
	fmt.Println("Encoded gob")

	if err := encoder.Encode(frame); err != nil {
		return err
	}

	fmt.Println("SENT META DATA")
	return nil
}

type MetaDataClient struct {
	Client http.Client
}

//Meta data server functions
func (c *MetaDataClient) getFileRecord(username string, password string, key string) (*records.Record, error) {

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

func (c *MetaDataClient) deleteFileRecord(username string, password string, key string) error {

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

func (c *MetaDataClient) createFileRecord(username string, password string, key string, fileName string, fileSize int64) error {

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

func (c *MetaDataClient) addAllowedUser(username string, password string, key string, allowedUser string) error {
	url := fmt.Sprintf("http://localhost:8080/record/allowedUsers?allowedUser=%s&username=%s&password=%s&key=%s",
		allowedUser,
		username,
		password,
		key)

	request, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(request)
	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error has occured while giving user access to record: %s\n", string(body))
	}

	return nil
}

func (c *MetaDataClient) removeAllowedUser(username string, password string, key string, removedUser string) error {
	url := fmt.Sprintf("http://localhost:8080/record/allowedUsers?removedUser=%s&username=%s&password=%s&key=%s",
		removedUser,
		username,
		password,
		key)

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(request)
	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error has occured while removing user: %s\n", string(body))
	}

	return nil
}
