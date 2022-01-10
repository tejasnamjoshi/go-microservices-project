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

func (i *Infrastructure) InitRouter(c *controllers.Auth, logger logging.Logger) {
	httpRouter := chi.NewRouter()

	// Initialize the Middleware
	httpRouter.Use(middleware.Logger)
	httpRouter.Use(i.Response)
	httpRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost*"},
		AllowedMethods: []string{"GET", "POST", "PATCH"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	// Initialize the routes
	httpRouter.Post("/user", c.AddUser)
	httpRouter.Get("/user/authorized", c.DecodeToken)
	httpRouter.Get("/login", c.Login)

	httpRouter.With(i.IsAuthorized).Get("/users", c.GetUsers)
	httpRouter.With(i.IsAuthorized).Delete("/user/{username}", c.DeleteUser)

	// Start the listener
	port := os.Getenv("AUTH_PORT")
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), httpRouter)
	if err != nil {
		logger.Fatal("Error starting server: %v", err.Error())
	}
}
