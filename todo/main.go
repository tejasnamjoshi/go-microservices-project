package main

import (
	"go-todo/todo/data"
	"go-todo/todo/handlers"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	l := log.New(os.Stdout, "go-todo", log.LstdFlags)
	l.Println("Welcome to the TODOS App")
	h := handlers.NewTodos(l)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.With(data.IsAuthorized).Get("/todos", h.GetByUsername)
	r.With(data.IsAuthorized).Post("/todo", h.CreateNewTodo)
	r.Patch("/todo/completed/{todoId}", h.MarkAsComplete)

	err := http.ListenAndServe(":3002", r)
	if err != nil {
		l.Fatal(err)
	}
}
