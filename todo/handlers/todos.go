package handlers

import (
	"go.uber.org/zap"
)

type Todos struct {
	l *zap.SugaredLogger
}

func NewTodos(l *zap.SugaredLogger) *Todos {
	return &Todos{l}
}
