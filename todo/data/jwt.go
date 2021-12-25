package data

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
	UserId     int64  `json:"userId"`
	Username   string `json:"username"`
	Authorized bool   `json:"authorized"`
}

func getNats() (*nats.Conn, error) {
	var nc *nats.Conn
	uri := os.Getenv("NATS_URI")
	var err error

	for i := 0; i < 5; i++ {
		nc, err = nats.Connect(uri)
		if err == nil {
			break
		}

		fmt.Println("Waiting before connecting to NATS at:", uri)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("Error establishing connection to NATS: %s", err)
	}
	fmt.Println("Connected to NATS at:", nc.ConnectedUrl())

	return nc, nil
}

func IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized user access."))
			return
		}
		nc, err := getNats()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Error connecting to NATS."))
			return
		}

		msg, err := nc.Request("authenticate", []byte(authHeader), time.Minute)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized user access."))
			return
		}
		var userClaims = &UserClaims{}
		err = json.Unmarshal(msg.Data, userClaims)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Error parsing token."))
			return
		}
		if !userClaims.Authorized {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized user access."))
			return
		}
		ctx := context.WithValue(r.Context(), "userId", userClaims.UserId)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
