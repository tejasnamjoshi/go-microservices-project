package main

import (
	"fmt"
	"go-todo/todo/handlers"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	l := getLogger()
	err := godotenv.Load(".env")
	if err != nil {
		l.Error("Error loading .env file")
	}
	port := os.Getenv("TODO_PORT")

	h := handlers.NewTodos(l)

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

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		l.Fatal(err)
	}

	l.Info("Welcome to the TODOS App")
}

func getLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	return sugar

	// l := log.New(os.Stdout, "go-todo", log.LstdFlags)
}
