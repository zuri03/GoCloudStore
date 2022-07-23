package records

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/zuri03/GoCloudStore/common"
)

type mockUserDB struct{}

func (mock mockUserDB) GetUser(id string) (*common.User, error) {
	return nil, nil
}

func (mock mockUserDB) CreateUser(user *common.User) error {
	return nil
}

func (mock mockUserDB) SearchUser(username, password string) ([]*common.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error in bcrypt: %s\n", err.Error())
		return nil, nil
	}
	return []*common.User{
		{
			Id:       "id",
			Username: "user",
			Password: hashedPassword,
		},
	}, nil
}

func TestServeHTTP(t *testing.T) {
	type scenario struct {
		Name               string
		username           string
		password           string
		ExpectedStatusCode int
		ExpectedResult     string
		Middleware         func(username, password string) error
	}

	scenarios := []scenario{
		{
			Name:               "Empty Credentials Tests",
			username:           "",
			password:           "",
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResult:     "password, username missing from request\n",
			Middleware: func(username, password string) error {
				return fmt.Errorf("password, username missing from request")
			},
		},
		{
			Name:               "Empty Password Test",
			username:           "user",
			password:           "",
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResult:     "username missing from request\n",
			Middleware: func(username, password string) error {
				return fmt.Errorf("username missing from request")
			},
		},
		{
			Name:               "Success Test",
			username:           "user",
			password:           "pass",
			ExpectedResult:     "{\"id\":\"id\"}",
			ExpectedStatusCode: http.StatusOK,
			Middleware: func(username, password string) error {
				return nil
			},
		},
	}

	baseUrl := "http://localhost:8080/auth?"
	testWaitGroup := new(sync.WaitGroup)
	for _, scen := range scenarios {
		url := fmt.Sprintf("%susername=%s&password=%s", baseUrl, scen.username, scen.password)
		testRequest, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Errorf("Error making new request: %s\n", err.Error())
			return
		}

		t.Run(scen.Name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()

			handler := &AuthHandler{
				dbClient:       mockUserDB{},
				routineTracker: testWaitGroup,
				validateParams: scen.Middleware,
			}

			handler.ServeHTTP(responseRecorder, testRequest)
			responseResult := responseRecorder.Result()

			if responseResult.StatusCode != scen.ExpectedStatusCode {
				t.Errorf("incorrect status code excepted %d got %d \n", scen.ExpectedStatusCode, responseResult.StatusCode)
			}

			if responseRecorder.Body.String() != scen.ExpectedResult {
				t.Errorf("incorrect return expected \"%s\" got \"%s\" \n", scen.ExpectedResult, responseRecorder.Body.String())
			}
		})
	}
}
