package infrastructure

import (
	"encoding/json"
	"go-todo/auth/service"
	"net/http"
	"time"
)

// A middleware that sets the content type on every request
func (m *Infrastructure) Response(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// Authorization middleware that returns whether a user has a valid JWT or not
func (m *Infrastructure) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		nc, err := m.GetNats()
		if err != nil {
			m.Logger.Error(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		authHeader := r.Header.Get("Authorization")
		msg, err := nc.Request("authenticate", []byte(authHeader), time.Minute)
		if err != nil {
			m.Logger.Error(err.Error())
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		var userClaims = &service.CustomClaims{}
		err = json.Unmarshal(msg.Data, userClaims)
		if err != nil {
			m.Logger.Error("Cannot authenticate")
			return
		}
		if !userClaims.Authorized {
			m.Logger.Error("Cannot authenticate")
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		m.Logger.Info("Authenticated")
		next.ServeHTTP(rw, r)
	})
}
