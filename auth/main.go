package main

import (
	"go-todo/auth/controllers"
	"go-todo/auth/infrastructure"
	"go-todo/auth/logging"
	"go-todo/auth/repository"
	"go-todo/auth/response"
	"go-todo/auth/service"
	"go-todo/auth/store"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	logger := logging.NewZapLogger()
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("Cannot load .env")
	}

	db := store.GetDb(logger)
	userRepository := repository.NewMysqlRepository(db, logger)
	jwtService := service.NewJWTService(logger)
	resp := response.NewResponse(logger)
	c := controllers.NewAuthController(controllers.App{
		Validator:   validator.New(),
		Logger:      logger,
		UserService: service.NewUserService(userRepository, logger, jwtService),
		JwtService:  jwtService,
		Response:    resp,
	})

	i := infrastructure.NewInfrastructure(logger, jwtService)
	i.InitNats()
	logger.Info("Welcome to the AUTH App")
	i.InitRouter(c, logger)
}
