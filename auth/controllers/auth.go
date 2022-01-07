package controllers

import (
	"go-todo/auth/logging"
	"go-todo/auth/service"

	"github.com/go-playground/validator/v10"
)

type App struct {
	Validator   *validator.Validate
	Logger      logging.Logger
	JwtService  service.JWTService
	UserService service.UserService
}

type Auth struct {
	App
}

func NewAuthController(a App) *Auth {
	return &Auth{a}
}
