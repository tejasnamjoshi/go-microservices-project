package data

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type UserClaims struct {
	UserId     int64  `json:"userId"`
	Username   string `json:"username"`
	Authorized bool   `json:"authorized"`
}

func IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized user access."))
			return
		}
		req, err := http.NewRequest(http.MethodGet, "http://localhost:3001/user/authorized", nil)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		q := req.URL.Query()
		q.Add("token", strings.Split(authHeader, " ")[1])
		req.URL.RawQuery = q.Encode()
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		if resp.StatusCode != http.StatusOK {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		userClaims := &UserClaims{}
		defer resp.Body.Close()
		d := json.NewDecoder(resp.Body)
		d.Decode(userClaims)

		ctx := context.WithValue(r.Context(), "userId", userClaims.UserId)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
