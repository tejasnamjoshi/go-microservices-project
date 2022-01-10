package entities

import (
	"encoding/json"
	"io"
)

type User struct {
	Id       int    `db:"id"`
	Username string `db:"username" validate:"required"`
	Password string `db:"password" validate:"passwd"`
}

type Users []User

// Decodes from io reader to the user object.
func (user *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(user)
}
