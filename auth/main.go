package main

import (
	"log"
	"net/http"
	"os"

	"go-todo/auth/data"
	"go-todo/auth/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	l := log.New(os.Stdout, "go-todo", log.LstdFlags)

	l.Println("Welcome to the AUTH App")
	h := handlers.NewAuth(l)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.With(data.IsAuthorized).Get("/users", h.GetUsers)
	r.With(data.IsAuthorized).Post("/user", h.AddUser)
	r.With(data.IsAuthorized).Delete("/user/{username}", h.DeleteUser)
	r.Get("/user/authorized", h.GetUserAuthStatus)

	r.Post("/login", h.Login)

	err = http.ListenAndServe(":3001", r)
	if err != nil {
		l.Fatal(err)
	}
}
