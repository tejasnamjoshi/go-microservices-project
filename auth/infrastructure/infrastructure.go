package infrastructure

import (
	"go-todo/auth/logging"
	"go-todo/auth/service"
)

type Infrastructure struct {
	Logger     logging.Logger
	JwtService service.JWTService
}

func NewInfrastructure(logger logging.Logger, jwtService service.JWTService) *Infrastructure {
	return &Infrastructure{Logger: logger, JwtService: jwtService}
}
