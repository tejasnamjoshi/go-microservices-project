package handlers

import (
	"go.uber.org/zap"
)

type Auth struct {
	l *zap.SugaredLogger
}

func NewAuth(l *zap.SugaredLogger) *Auth {
	return &Auth{l}
}
