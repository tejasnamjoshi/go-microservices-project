package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"go-todo/auth/data"
	"go-todo/auth/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
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

	r.With(data.IsAuthorized).Get("/users", h.GetUsers)
	r.With(data.IsAuthorized).Post("/user", h.AddUser)
	r.With(data.IsAuthorized).Delete("/user/{username}", h.DeleteUser)
	r.Get("/user/authorized", h.GetUserAuthStatus)

	r.Post("/login", h.Login)

	InitNats(l)

	err = http.ListenAndServe(":3001", r)
	if err != nil {
		l.Fatal(err)
	}
}

func InitNats(l *log.Logger) {
	nc, err := data.GetNats()
	if err != nil {
		l.Fatal(err)
	}

	fmt.Println("Listening for auth", nc.ConnectedUrl())
	nc.Subscribe("authenticate", func(msg *nats.Msg) {
		authHeader := string(msg.Data)
		jwtParts := strings.Split(authHeader, " ")
		if len(jwtParts) <= 1 {
			nc.Publish(msg.Reply, []byte("No token available"))
			return
		}

		var claims = &data.CustomClaims{}
		ok, err := data.GetAuthorizationStatus(authHeader)
		if ok {
			claims, err = data.ParseJWT(jwtParts[1])
			if err != nil {
				fmt.Println("Error parsing token")
			}
		}
		resp, err := json.Marshal(*claims)
		nc.Publish(msg.Reply, []byte(resp))
	})
}
