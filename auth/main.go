package main

import (
	"log"
	"net/http"
	"os"

	"go-todo/auth/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	l := log.New(os.Stdout, "go-todo", log.LstdFlags)

	_ = hclog.New(&hclog.LoggerOptions{
		Name:  "my-app",
		Level: hclog.LevelFromString("DEBUG"),
	})

	l.Println("Welcome to the AUTH App")
	h := handlers.NewAuth(l)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.With(h.IsAuthorized).Get("/users", h.GetUsers)
	r.With(h.IsAuthorized).Post("/user", h.AddUser)
	r.With(h.IsAuthorized).Delete("/user/{username}", h.DeleteUser)
	r.Get("/user/authorized", h.GetUserAuthStatus)

	r.Post("/login", h.Login)

	h.InitNats()

	err = http.ListenAndServe(":3001", r)
	if err != nil {
		l.Fatal(err)
	}
}
