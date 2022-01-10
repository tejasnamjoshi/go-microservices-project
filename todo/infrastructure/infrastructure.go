package infrastructure

import (
	"go-todo/todo/controllers"
	"go-todo/todo/logging"
)

type Infrastructure struct {
	Logger     logging.Logger
	Controller *controllers.Todos
}

// Constructor function for the infrastructure package
func NewInfrastructure(logger logging.Logger, controller *controllers.Todos) *Infrastructure {
	return &Infrastructure{Logger: logger, Controller: controller}
}
