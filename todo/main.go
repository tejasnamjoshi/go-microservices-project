package main

import (
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

	r.Route("/todos/{username}", func(r chi.Router) {
		r.Get("/", h.GetByUsername)
		// r.Put("/", h.GetAll)
	})
	r.Post("/todo/{userId}", h.CreateNewTodo)

	err := http.ListenAndServe(":3002", r)
	if err != nil {
		l.Fatal(err)
	}
}
