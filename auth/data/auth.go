package data

import (
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int64  `db:"id"`
	Username string `db:"username" validate:"required"`
	Password string `db:"password" validate:"passwd"`
}

type Users []User

func (users Users) ToJSON(rw http.ResponseWriter) error {
	e := json.NewEncoder(rw)
	return e.Encode(users)
}

func (user *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(user)
}

func (user *User) GeneratePassword() error {
	pwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	user.Password = string(pwd)
	return nil
}

func (user *User) ComparePassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return true
	}

	return false
}
