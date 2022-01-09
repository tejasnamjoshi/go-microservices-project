package infrastructure

import (
	"context"
	"encoding/json"
	"go-todo/todo/controllers"
	"net/http"
	"time"
)

type UserClaims struct {
	UserId     int    `json:"userId"`
	Username   string `json:"username"`
	Authorized bool   `json:"authorized"`
}

func (m *Infrastructure) ResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (m *Infrastructure) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.Logger.Error("Unauthorized user access.")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized user access."))
			return
		}
		nc, err := m.GetNats()
		if err != nil {
			m.Logger.Error("Error connecting to NATS.")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Error connecting to NATS."))
			return
		}
		msg, err := nc.Request("authenticate", []byte(authHeader), time.Minute)
		if err != nil {
			m.Logger.Error("Unauthorized user access.")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized user access."))
			return
		}
		var userClaims = &UserClaims{}
		err = json.Unmarshal(msg.Data, userClaims)
		if err != nil {
			m.Logger.Error("Error parsing token.")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Error parsing token."))
			return
		}
		if !userClaims.Authorized {
			m.Logger.Error("unauthorized user access.")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("unauthorized user access."))
			return
		}

		ctx := context.WithValue(context.Background(), controllers.UserIdContext{}, userClaims.UserId)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
