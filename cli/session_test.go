package cli

import (
	"testing"
)

type scenario struct {
	Name             string
	Args             []string
	ExpectedUsername string
	ExpectedPassword string
	ExpectedId       string
}

//Testing the Parse args function
//TODO: Add command line tests
func TestProperSessionCreation(t *testing.T) {

	scenarios := []scenario{
		{
			Name:             "Testing with simple and proper arguments",
			Args:             []string{"username", "password"},
			ExpectedUsername: "username",
			ExpectedPassword: "password",
			ExpectedId:       "",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			testSession := ParseArgsIntoSession(scenario.Args)

			if testSession.Username != scenario.ExpectedUsername {
				t.Errorf("Expected username '%s', got %s\n", scenario.ExpectedUsername, testSession.Username)
			}

			if testSession.Password != scenario.ExpectedPassword {
				t.Errorf("Expected password '%s', got %s\n", scenario.ExpectedPassword, testSession.Password)
			}

			if testSession.Id != scenario.ExpectedId {
				t.Errorf("Excepted id %s, got %s\n", scenario.ExpectedId, testSession.Id)
			}
		})
	}
}
