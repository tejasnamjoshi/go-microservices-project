package data

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type CustomClaims struct {
	UserId     int64  `json:"userId"`
	Username   string `json:"username"`
	Authorized bool   `json:"authorized"`
	jwt.StandardClaims
}

func (user *User) GetJWT() (string, error) {
	var mySigningKey = []byte(os.Getenv("SECRET_KEY"))

	claims := CustomClaims{
		user.Id,
		user.Username,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
			Issuer:    os.Getenv("JWT_ISS"),
			Audience:  os.Getenv("JWT_AUD"),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return claims, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return claims, fmt.Errorf("Error decoding token")
	}
	return claims, nil
}

func IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		ok, err := GetAuthorizationStatus(authHeader)
		if err != nil {
			fmt.Println(err)
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		if ok {
			next.ServeHTTP(rw, r)
		}
	})
}

func GetAuthorizationStatus(authHeader string) (bool, error) {
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
