package handlers

import "log"

type Todos struct {
	l *log.Logger
}

func NewTodos(l *log.Logger) *Todos {
	return &Todos{l}
}
