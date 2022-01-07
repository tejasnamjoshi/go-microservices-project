package handlers

import (
	"go-todo/todo/logging"
	"go-todo/todo/repository"

	"github.com/jmoiron/sqlx"
)

type Todos struct {
	l              logging.Logger
	db             *sqlx.DB
	todoRepository repository.TodoRepository
}

func NewTodos(l logging.Logger, db *sqlx.DB, todoRepository repository.TodoRepository) *Todos {
	return &Todos{l, db, todoRepository}
}
