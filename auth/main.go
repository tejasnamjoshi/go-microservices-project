package main

import (
	"log"
	"net/http"
	"os"

	"go-todo/auth/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	l := log.New(os.Stdout, "go-todo", log.LstdFlags)
	l.Println("Welcome to the AUTH App")
	h := handlers.NewAuth(l)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/users", h.GetUsers)
	r.Post("/user", h.AddUser)
	r.Delete("/user/{username}", h.DeleteUser)

	err := http.ListenAndServe(":3001", r)
	if err != nil {
		l.Fatal(err)
	}
}
