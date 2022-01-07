package main

import (
	"fmt"
	"go-todo/todo/controllers"
	"go-todo/todo/logging"
	"go-todo/todo/repository"
	"go-todo/todo/service"
	"go-todo/todo/store"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	logger := logging.NewZapLogger()
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("Error loading .env file")
	}
	port := os.Getenv("TODO_PORT")

	db := store.GetDb(logger)
	todoRepository := repository.NewMysqlRepository(db, logger)
	todoService := service.NewTodoService(todoRepository, logger)
	h := controllers.NewTodos(controllers.App{
		todoService,
		logger,
	})

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost*"},
		AllowedMethods: []string{"GET", "POST", "PATCH"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))
	r.With(h.IsAuthorized).Get("/todos", h.GetByUsername)
	r.With(h.IsAuthorized).Post("/todos", h.CreateNewTodo)
	r.With(h.IsAuthorized).Patch("/todos/{todoId}", h.MarkAsComplete)

	logger.Info("Welcome to the TODOS App")
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		logger.Fatal(err.Error())
	}
}
