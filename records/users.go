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

type Users struct {
	userList map[string]User
}

func NewUsers() *Users {
	return &Users{
		userList: make(map[string]User),
	}
}

func (u *Users) create(username string, password string) (*User, error) {
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
	key := fmt.Sprintf("%s:%s", username, string(hash))
	u.userList[key] = user
	return &user, nil
}

func (u *Users) exists(username string, password string) (bool, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return true, err
	}
	key := fmt.Sprintf("%s:%s", username, hash)
	fmt.Printf("Searching exist with key: %s\n", key)
	_, ok := u.userList[key]
	return ok, nil
}

func (u *Users) get(username string, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("%s:%s", username, hash)
	fmt.Printf("Searching with key: %s\n", key)
	user, ok := u.userList[key]
	if !ok {
		return nil, nil
	}

	return &user, nil
}
