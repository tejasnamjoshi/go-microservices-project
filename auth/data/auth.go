package data

import (
	"encoding/json"
	"io"
	"net/http"
)

type User struct {
	Id       int64  `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
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
