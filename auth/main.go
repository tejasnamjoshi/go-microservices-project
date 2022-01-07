package main

import (
	"os"

	"go-todo/auth/controllers"
	"go-todo/auth/logging"
	"go-todo/auth/repository"
	"go-todo/auth/router"
	"go-todo/auth/service"
	"go-todo/auth/store"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	logger := logging.NewZapLogger()
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("Cannot load .env")
	}

	Validate := validator.New()
	db := store.GetDb(logger)

	userRepository := repository.NewMysqlRepository(db, logger)
	userService := service.NewUserService(userRepository, logger)
	jwtService := service.NewJWTService(logger)

	c := controllers.NewAuthController(controllers.App{
		Validator:   Validate,
		Logger:      logger,
		UserService: userService,
		JwtService:  jwtService,
	})
	logger.Info("Welcome to the AUTH App")
	initRouter(c)
}

func initRouter(c *controllers.Auth) {
	httpRouter := router.NewChiRouter()
	httpRouter.USE(middleware.Logger)
	httpRouter.USE(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost*"},
		AllowedMethods: []string{"GET", "POST", "PATCH"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	c.InitNats()
	httpRouter.POST("/user", c.AddUser)
	httpRouter.GET("/user/authorized", c.DecodeToken)
	httpRouter.GET("/login", c.Login)

	httpRouter.WITH(c.IsAuthorized).GET("/users", c.GetUsers)
	httpRouter.WITH(c.IsAuthorized).DELETE("/user/{username}", c.DeleteUser)

	port := os.Getenv("AUTH_PORT")
	httpRouter.SERVE(port)
}
