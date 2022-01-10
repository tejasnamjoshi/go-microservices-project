package entities

import (
	"encoding/json"
	"io"
)

type Todo struct {
	Id        int    `db:"id"`
	Content   string `db:"content"`
	Completed bool   `db:"completed"`
}

type Todos []Todo

func (todo *Todo) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(todo)
}
