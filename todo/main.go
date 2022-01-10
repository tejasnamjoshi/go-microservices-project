package main

import (
	"go-todo/todo/controllers"
	"go-todo/todo/infrastructure"
	"go-todo/todo/logging"
	"go-todo/todo/repository"
	"go-todo/todo/response"
	"go-todo/todo/service"
	"go-todo/todo/store"

	"github.com/joho/godotenv"
)

func main() {
	// Creates logger
	logger := logging.NewZapLogger()
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("Error loading .env file")
	}

	// Creates DB and initializes controllers
	db := store.GetDb(logger)
	todoRepository := repository.NewMysqlRepository(db, logger)
	todoService := service.NewTodoService(todoRepository, logger)
	resp := response.NewResponse(logger)
	c := controllers.NewTodos(controllers.App{
		TodoService: todoService,
		Logger:      logger,
		Response:    resp,
	})

	// Sets up routing, middleware, NATS
	i := infrastructure.NewInfrastructure(logger, c)
	i.InitRouter()
}
