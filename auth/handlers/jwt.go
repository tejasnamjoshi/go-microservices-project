package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-todo/auth/data"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
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

func (a Auth) GetAuthorizationStatus(authHeader string) (bool, error) {
	var mySigningKey = []byte(os.Getenv("SECRET_KEY"))
	jwtParts := strings.Split(authHeader, " ")
	if len(jwtParts) <= 1 {
		return false, errors.New("No Authorization Token provided")
	}
	token, err := jwt.Parse(jwtParts[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf(("Invalid Signing Method"))
		}
		aud := os.Getenv("JWT_AUD")
		checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
		if !checkAudience {
			return false, fmt.Errorf(("invalid aud"))
		}
		iss := os.Getenv("JWT_ISS")
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
		if !checkIss {
			return false, fmt.Errorf(("invalid iss"))
		}

		return mySigningKey, nil
	})
	if err != nil {
		return false, err
	}

	if token.Valid {
		return true, nil
	}

	return false, errors.New("Invalid Token")
}
