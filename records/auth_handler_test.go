package records

import (
	"fmt"
	"net/http"
	"testing"
)

type scenario struct {
	username           string
	password           string
	ExpectedStatusCode int
}

func TestServeHTTP(t *testing.T) {

	mockMiddleware := func(username, password string) bool {
		return false
	}

	scenarios := []scenario{
		scenario{
			username:           "",
			password:           "",
			ExpectedStatusCode: http.StatusBadRequest,
		},
		scenario{
			username:           "user",
			password:           "",
			ExpectedStatusCode: http.StatusBadRequest,
		},
		scenario{
			username:           "user",
			password:           "pass",
			ExpectedStatusCode: http.StatusBadRequest,
		},
	}

	requests := []*http.Request{}
	baseUrl := "http://localhost:8080/auth?"
	for idx, scen := range scenarios {
		url := fmt.Sprintf("%susername=%s&password=%s", baseUrl, scen.username, scen.password)
		testRequest, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Error making reqs: %s\n", err.Error())
			return
		}
		requests = append(requests, testRequest)
	}

}
