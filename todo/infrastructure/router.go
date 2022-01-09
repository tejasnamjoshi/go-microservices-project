package infrastructure

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (i Infrastructure) InitRouter() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost*"},
		AllowedMethods: []string{"GET", "POST", "PATCH"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))
	r.With(i.IsAuthorized).Get("/todos", i.Controller.GetByUsername)
	r.With(i.IsAuthorized).Post("/todos", i.Controller.CreateNewTodo)
	r.With(i.IsAuthorized).Patch("/todos/{todoId}", i.Controller.MarkAsComplete)

	i.Logger.Info("Welcome to the TODOS App")
	port := os.Getenv("TODO_PORT")
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		i.Logger.Fatal(err.Error())
	}
}
