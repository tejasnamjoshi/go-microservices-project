package infrastructure

import (
	"fmt"
	"go-todo/todo/controllers"
	"go-todo/todo/logging"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func InitRouter(c *controllers.Todos, logger logging.Logger) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost*"},
		AllowedMethods: []string{"GET", "POST", "PATCH"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))
	r.With(c.IsAuthorized).Get("/todos", c.GetByUsername)
	r.With(c.IsAuthorized).Post("/todos", c.CreateNewTodo)
	r.With(c.IsAuthorized).Patch("/todos/{todoId}", c.MarkAsComplete)

	logger.Info("Welcome to the TODOS App")
	port := os.Getenv("TODO_PORT")
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		logger.Fatal(err.Error())
	}
}
