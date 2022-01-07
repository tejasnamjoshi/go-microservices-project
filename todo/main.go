package main

import (
	"go-todo/todo/controllers"
	"go-todo/todo/infrastructure"
	"go-todo/todo/logging"
	"go-todo/todo/repository"
	"go-todo/todo/service"
	"go-todo/todo/store"

	"github.com/joho/godotenv"
)

func main() {
	logger := logging.NewZapLogger()
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("Error loading .env file")
	}

	db := store.GetDb(logger)
	todoRepository := repository.NewMysqlRepository(db, logger)
	todoService := service.NewTodoService(todoRepository, logger)
	c := controllers.NewTodos(controllers.App{
		TodoService: todoService,
		Logger:      logger,
	})

	infrastructure.InitRouter(c, logger)
}
