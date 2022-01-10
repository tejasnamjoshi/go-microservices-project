package infrastructure

import (
	"go-todo/auth/logging"
	"go-todo/auth/service"
)

type Infrastructure struct {
	Logger     logging.Logger
	JwtService service.JWTService
}

// Constructor function for the infrastructure package
func NewInfrastructure(logger logging.Logger, jwtService service.JWTService) *Infrastructure {
	return &Infrastructure{Logger: logger, JwtService: jwtService}
}
