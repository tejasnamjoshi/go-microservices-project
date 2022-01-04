package main

import (
	"os"

	"go-todo/auth/controllers"
	"go-todo/auth/repository"
	"go-todo/auth/router"
	"go-todo/auth/service"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	Validate       *validator.Validate       = validator.New()
	httpRouter     router.Router             = router.NewChiRouter()
	userRepository repository.UserRepository = repository.NewMysqlRepository()
	userService    service.UserService       = service.NewUserService(userRepository)
	jwtService     service.JWTService        = service.NewJWTService()
)

func main() {
	l := getLogger()
	Validate.RegisterValidation("passwd", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > 5
	})

	err := godotenv.Load(".env")
	if err != nil {
		l.Error("Cannot load .env")
	}
	port := os.Getenv("AUTH_PORT")

	c := controllers.NewAuthController(controllers.App{
		Validator:   Validate,
		Logger:      getLogger(),
		UserService: userService,
		JwtService:  jwtService,
	})

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

	httpRouter.SERVE(port)

	l.Info("Welcome to the AUTH App")
}

func getLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	return sugar
}
