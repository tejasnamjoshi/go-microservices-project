package data

import (
	"fmt"
	"os"
	"strconv"
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
	st, _ := strconv.Atoi(os.Getenv("SESSION_TIME"))
	claims := CustomClaims{
		user.Id,
		user.Username,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(st)).Unix(),
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
