package entities

import (
	"encoding/json"
	"io"
	"net/http"
)

type Todo struct {
	Id        int    `db:"id"`
	Content   string `db:"content"`
	Completed bool   `db:"completed"`
}

type Todos []Todo

func (todos Todos) ToJSON(rw http.ResponseWriter) error {
	e := json.NewEncoder(rw)
	return e.Encode(todos)
}

func (todo *Todo) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(todo)
}
