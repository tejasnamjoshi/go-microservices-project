package controllers

import (
	"encoding/json"
	"go-todo/auth/data"
	"go-todo/auth/service"
	"net/http"
	"time"
)

func (a Auth) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		nc, err := data.GetNats(a.Logger)
		if err != nil {
			a.Logger.Error(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		authHeader := r.Header.Get("Authorization")
		msg, err := nc.Request("authenticate", []byte(authHeader), time.Minute)
		if err != nil {
			a.Logger.Error(err.Error())
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		var userClaims = &service.CustomClaims{}
		err = json.Unmarshal(msg.Data, userClaims)
		if err != nil {
			a.Logger.Error("Cannot authenticate")
			return
		}
		if !userClaims.Authorized {
			a.Logger.Error("Cannot authenticate")
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		a.Logger.Info("Authenticated")
		next.ServeHTTP(rw, r)
	})
}
