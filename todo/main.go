package main

import (
	"fmt"
	"go-todo/todo/handlers"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	port := os.Getenv("TODO_PORT")

	l := log.New(os.Stdout, "go-todo", log.LstdFlags)
	l.Println("Welcome to the TODOS App")
	h := handlers.NewTodos(l)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.With(h.IsAuthorized).Get("/todos", h.GetByUsername)
	r.With(h.IsAuthorized).Post("/todo", h.CreateNewTodo)
	r.Patch("/todo/completed/{todoId}", h.MarkAsComplete)

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		l.Fatal(err)
	}
}
