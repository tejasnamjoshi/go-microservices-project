package infrastructure

import (
	"fmt"
	"go-todo/auth/controllers"
	"go-todo/auth/logging"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func InitRouter(c *controllers.Auth, logger logging.Logger) {
	httpRouter := chi.NewRouter()
	httpRouter.Use(middleware.Logger)
	httpRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost*"},
		AllowedMethods: []string{"GET", "POST", "PATCH"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	httpRouter.Post("/user", c.AddUser)
	httpRouter.Get("/user/authorized", c.DecodeToken)
	httpRouter.Get("/login", c.Login)

	httpRouter.With(c.IsAuthorized).Get("/users", c.GetUsers)
	httpRouter.With(c.IsAuthorized).Delete("/user/{username}", c.DeleteUser)

	port := os.Getenv("AUTH_PORT")
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), httpRouter)
	if err != nil {
		logger.Fatal("Error starting server: %v", err.Error())
	}
}
