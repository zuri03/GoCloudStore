package records

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	Password     []byte `json:"password"`
	CreationDate string `json:"creationDate"`
}

//For now the username will serve as the key until persistant storage is added
type Users struct {
	userList map[string]User
}

func NewUsers() *Users {
	return &Users{
		userList: make(map[string]User),
	}
}

func (u *Users) new(username string, password string) (*User, error) {

	if _, ok := u.userList[username]; ok {
		return nil, fmt.Errorf("User %s already exists\n", username)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	now := time.Now().Format("2006-01-02 03:04:05")
	id := uuid.New()

	user := User{
		Id:           id.String(),
		Username:     username,
		Password:     hash,
		CreationDate: now,
	}
	fmt.Printf("Created user => %s\n", username)
	u.userList[username] = user
	return &user, nil
}

func (u *Users) exists(username string, password string) (bool, error) {
	fmt.Printf("Searching exist with key: %s\n", username)
	_, ok := u.userList[username]
	return ok, nil
}

func (u *Users) get(username string, password string) (*User, error) {
	fmt.Printf("Searching with key: %s\n", username)
	user, ok := u.userList[username]
	if !ok {
		return nil, nil
	}

	return &user, nil
}
