package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

type UserClaims struct {
	UserId     int    `json:"userId"`
	Username   string `json:"username"`
	Authorized bool   `json:"authorized"`
}

func (t Todos) GetNats() (*nats.Conn, error) {
	var nc *nats.Conn
	uri := os.Getenv("NATS_URI")
	var err error

	for i := 0; i < 5; i++ {
		nc, err = nats.Connect(uri)
		if err == nil {
			break
		}

		t.Logger.Warn("Waiting before connecting to NATS at:", uri)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("error establishing connection to NATS: %s", err)
	}
	t.Logger.Info("Connected to NATS at:", nc.ConnectedUrl())

	return nc, nil
}

func (t Todos) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			t.Logger.Error("Unauthorized user access.")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized user access."))
			return
		}
		nc, err := t.GetNats()
		if err != nil {
			t.Logger.Error("Error connecting to NATS.")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Error connecting to NATS."))
			return
		}
		msg, err := nc.Request("authenticate", []byte(authHeader), time.Minute)
		if err != nil {
			t.Logger.Error("Unauthorized user access.")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized user access."))
			return
		}
		var userClaims = &UserClaims{}
		err = json.Unmarshal(msg.Data, userClaims)
		if err != nil {
			t.Logger.Error("Error parsing token.")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Error parsing token."))
			return
		}
		if !userClaims.Authorized {
			t.Logger.Error("unauthorized user access.")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("unauthorized user access."))
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userClaims.UserId)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
