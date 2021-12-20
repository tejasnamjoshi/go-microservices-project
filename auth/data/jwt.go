package data

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func (user *User) GetJWT() (string, error) {
	var mySigningKey = []byte(os.Getenv("SECRET_KEY"))
	fmt.Println("test", os.Getenv("SECRET_KEY"))
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = user.Username
	claims["aud"] = os.Getenv("JWT_AUD")
	claims["iss"] = os.Getenv("JWT_ISS")
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func IsAuthorized(next http.Handler) http.Handler {
	var mySigningKey = []byte(os.Getenv("SECRET_KEY"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] == nil {
			fmt.Fprintf(w, "No Authorization Token provided")
			return
		}

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(("Invalid Signing Method"))
			}
			aud := os.Getenv("JWT_AUD")
			checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAudience {
				return nil, fmt.Errorf(("invalid aud"))
			}
			// verify iss claim
			iss := os.Getenv("JWT_ISS")
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return nil, fmt.Errorf(("invalid iss"))
			}

			return mySigningKey, nil
		})
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		if token.Valid {
			next.ServeHTTP(w, r)
		}
	})
}
