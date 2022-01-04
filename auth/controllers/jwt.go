package controllers

import (
	"encoding/json"
	"fmt"
	"go-todo/auth/data"
	"net/http"
	"time"
)

func (a Auth) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		nc, err := data.GetNats()
		if err != nil {
			fmt.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		authHeader := r.Header.Get("Authorization")
		msg, err := nc.Request("authenticate", []byte(authHeader), time.Minute)
		if err != nil {
			fmt.Println(err)
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		var userClaims = &data.CustomClaims{}
		err = json.Unmarshal(msg.Data, userClaims)
		if err != nil {
			fmt.Println("Cannot authenticate")
			return
		}
		if !userClaims.Authorized {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(rw, r)
	})
}
