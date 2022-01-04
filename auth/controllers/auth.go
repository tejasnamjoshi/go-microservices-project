package controllers

import (
	"go-todo/auth/service"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type App struct {
	Validator   *validator.Validate
	Logger      *zap.SugaredLogger
	JwtService  service.JWTService
	UserService service.UserService
}

type Auth struct {
	App
}

func NewAuthController(a App) *Auth {
	return &Auth{a}
}
