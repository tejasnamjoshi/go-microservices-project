package controllers

import (
	"go-todo/auth/service"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Auth struct {
	l *zap.SugaredLogger
	v *validator.Validate
}

var (
	userService service.UserService = service.NewUserService()
	jwtService  service.JWTService  = service.NewJWTService()
)

func NewAuthController(l *zap.SugaredLogger, v *validator.Validate) *Auth {
	return &Auth{l, v}
}
