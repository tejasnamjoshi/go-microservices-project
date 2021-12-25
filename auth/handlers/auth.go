package handlers

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Auth struct {
	l *zap.SugaredLogger
	v *validator.Validate
}

func NewAuth(l *zap.SugaredLogger, v *validator.Validate) *Auth {
	return &Auth{l, v}
}
