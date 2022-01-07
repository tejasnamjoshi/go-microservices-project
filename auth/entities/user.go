package entities

import (
	"encoding/json"
	"io"
	"net/http"
)

type User struct {
	Id       int    `db:"id"`
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
