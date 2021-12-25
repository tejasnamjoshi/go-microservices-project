package main

import (
	"fmt"
	"net/http"
	"os"

	"go-todo/auth/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	l := getLogger()

	err := godotenv.Load(".env")
	if err != nil {
		l.Error("Cannot load .env")
	}
	port := os.Getenv("AUTH_PORT")

	h := handlers.NewAuth(l)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.With(h.IsAuthorized).Get("/users", h.GetUsers)
	r.Post("/user", h.AddUser)
	r.With(h.IsAuthorized).Delete("/user/{username}", h.DeleteUser)
	r.Get("/user/authorized", h.GetUserAuthStatus)

	r.Post("/login", h.Login)

	h.InitNats()

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	err = http.ListenAndServe(":3001", r)
	if err != nil {
		l.Fatal(err)
	}

	l.Info("Welcome to the AUTH App")
}

func getLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	return sugar

	// l := log.New(os.Stdout, "go-todo", log.LstdFlags)
}
