package controllers

import (
	"go-todo/todo/logging"
	"go-todo/todo/service"
)

type App struct {
	TodoService service.TodoService
	Logger      logging.Logger
}

type Todos struct {
	App
}

func NewTodos(a App) *Todos {
	return &Todos{a}
}
