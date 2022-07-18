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

	"github.com/zuri03/GoCloudStore/common"
)

type AuthenticationResponse struct {
	Id string `json:"id"`
}

func newEncoderDecorder(connection net.Conn) (*gob.Encoder, *gob.Decoder) {
	gob.Register(new(common.ProtocolFrame))
	gob.Register(new(common.FileMetaData))

	return gob.NewEncoder(connection), gob.NewDecoder(connection)
}

//This function forces the client to wait for the server to send a message letting
//the client know it is ready to move to the next step of the process
func waitForProceed(decoder *gob.Decoder) error {
	var frame common.ProtocolFrame
	if err := decoder.Decode(&frame); err != nil {
		return err
	}

	fmt.Printf("Proceed received => %+v\n", frame)

	if frame.Type == common.ERROR_FRAME {
		return fmt.Errorf("Error on server: %s\n", string(frame.Data))
	}

	if frame.Type != common.PROCEED_FRAME {
		return fmt.Errorf("Unexpected frame got: %d\n", frame.Type)
	}

	fmt.Printf("Got proceed frame => %d\n", frame.Type)
	return nil
}

func sendMetaDataToServer(frameType common.FrameType, meta common.FileMetaData, encoder *gob.Encoder) error {
	metaBuffer := new(bytes.Buffer)
	if err := gob.NewEncoder(metaBuffer).Encode(meta); err != nil {
		return err
	}
	fmt.Printf("encoding meta => %+v\n", meta)
	fmt.Printf("encoding meta length=> %d\n", metaBuffer.Len())
	frame := common.ProtocolFrame{
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
//If there is ever an error just return true so that the session hanlder does not assume the user does not exist
func (c *MetaDataClient) authenticate(username string, password string) (string, bool, error) {
	url := fmt.Sprintf("http://localhost:8080/auth?username=%s&password=%s", username, password)
	fmt.Printf("url => %s\n", url)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", false, err
	}

	resp, err := c.Client.Do(request)
	if err != nil {
		return "", false, err
	}

	fmt.Printf("Status code => %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return "", false, fmt.Errorf("Server returned error: %d\n", resp.StatusCode)
	}

	var id AuthenticationResponse
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", false, err
	}
	fmt.Printf("response => %s\n", string(responseBytes))
	if err = json.Unmarshal(responseBytes, &id); err != nil {
		return "", false, err
	}

	return id.Id, id.Id != "", nil
}

func (c *MetaDataClient) createUser(username string, password string) (*common.User, error) {

	url := fmt.Sprintf("http://localhost:8080/user?username=%s&password=%s", username, password)
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error has occured while creating user: %d\n", resp.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var user common.User
	if err := json.Unmarshal(responseBody, &user); err != nil {
		return nil, fmt.Errorf("Error unmarshaling response: %s\n", err.Error())
	}

	return &user, nil
}

func (c *MetaDataClient) getFileRecord(owner string, key string) (*common.Record, error) {

	url := fmt.Sprintf("http://localhost:8080/record?owner=%s&key=%s", owner, key)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(request)
	if err != nil {
		return nil, err
	}

	var record common.Record
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Server return error => %d\n", resp.StatusCode)
		if resp.StatusCode == http.StatusForbidden {
			return nil, fmt.Errorf("%s is not allowed to view %s\n", owner, key)
		}
		return nil, fmt.Errorf(string(body))
	}

	json.Unmarshal(body, &record)
	return &record, nil
}

func (c *MetaDataClient) deleteFileRecord(owner, key string) error {

	url := fmt.Sprintf("http://localhost:8080/record?owner=%s&key=%s", owner, key)
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("unauthorized")
	}

	var record common.Record
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

type Request struct {
	Owner    string `json:"owner"`
	Key      string `json:"key"`
	FileName string `json:"name"`
	Size     int    `json:"size"`
}

func (c *MetaDataClient) createFileRecord(owner, key, fileName string, fileSize int64) error {

	record := Request{
		Owner:    owner,
		Key:      key,
		FileName: fileName,
		Size:     int(fileSize),
	}
	//need to include necessary params to pass params check middleware
	url := fmt.Sprintf("http://localhost:8080/record?owner=%s&key=%s", owner, key)
	recordBytes, _ := json.Marshal(record)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(recordBytes))
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("unauthorized")
	}

	if resp.StatusCode != 200 {
		errorMessage, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error: record server returned error: %d: %s\n", resp.StatusCode, string(errorMessage))
	}

	return nil
}

func (c *MetaDataClient) addAllowedUser(owner, key, allowedUser string) error {
	url := fmt.Sprintf("http://localhost:8080/record/allowedUser?allowedUser=%s&owner=%s&key=%s", allowedUser, owner, key)

	request, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(request)
	if resp.StatusCode != http.StatusOK {

		if resp.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("unauthorized")
		}

		return fmt.Errorf("Error has occured while giving user access to record: %d\n", resp.StatusCode)
	}

	return nil
}

func (c *MetaDataClient) removeAllowedUser(owner, key string, removedUser string) error {
	url := fmt.Sprintf("http://localhost:8080/record/allowedUser?removedUser=%s&owner=%s&key=%s", removedUser, owner, key)

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(request)
	if resp.StatusCode != http.StatusOK {

		if resp.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("unauthorized")
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error has occured while removing user: %s\n", string(body))
	}

	return nil
}
