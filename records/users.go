package records

import (
	"time"

	_ "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	Password     []byte `json:"password"`
	CreationDate string `json:"creationDate"`
}

type Users map[string]User

func (u *Users) New() map[string]User {
	return make(map[string]User)
}

func (u *Users) create(username string, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	now := time.Now().Format("2006-01-02 03:04:05")
	user := User{
		Username:     username,
		Password:     hash,
		CreationDate: now,
	}
	return &user, nil
}

func (u *Users) exists(username string, password string) (bool, error) {
	return true, nil
}

func (u *Users) get(username string, password string) (*User, error) {
	return &User{}, nil
}
