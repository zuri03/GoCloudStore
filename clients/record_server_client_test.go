package clients

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

//More information about testing an Http Client https://engineering.teknasyon.com/how-to-write-unit-tests-for-http-clients-in-go-5de79aa4f92
type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// Test4xxError tests that a 4xx error is returned as an APIError.
func TestAuthenticateUser(t *testing.T) {
	hclient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) *http.Response {
			username := req.FormValue("username")
			password := req.FormValue("password")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader(fmt.Sprintf(`{"id":"%s:%s"}`, username, password))),
			}
		})}

	recordServerClient := RecordServerclient{HttpClient: *hclient}

	testUsername := "username"
	testPassword := "password"

	id, userExists, err := recordServerClient.AuthenticateUser(testUsername, testPassword)

	fmt.Printf("id %s, userExists %t \n", id, userExists)

	expectedId := fmt.Sprintf("%s:%s", testUsername, testPassword)

	if id != expectedId {
		t.Errorf("Incorrect Id, expected %s, got %s\n", expectedId, id)
	}

	if !userExists {
		t.Errorf("Expected userExists to be true but got false")
	}

	if err != nil {
		t.Errorf("Unexpected error: %s\n", err.Error())
	}

	expectedErrorMessage := "example error"
	expectedError := fmt.Sprintf("%d: %s\n", http.StatusForbidden, expectedErrorMessage)

	errorHttpClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusForbidden,
				Body:       ioutil.NopCloser(strings.NewReader(expectedErrorMessage)),
			}
		})}

	recordServerClient.HttpClient = *errorHttpClient

	id, userExists, err = recordServerClient.AuthenticateUser("", "")

	if id != "" {
		t.Errorf("Expected empty id string, got %s\n", id)
	}

	if userExists {
		t.Errorf("Expected user exists to be false but got true\n")
	}

	if err.Error() != expectedError {
		t.Errorf("Expected error incorrect %s\n", err.Error())
	}
}
