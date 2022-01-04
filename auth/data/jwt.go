package data

import (
	"github.com/golang-jwt/jwt"
)

type CustomClaims struct {
	UserId     int64  `json:"userId"`
	Username   string `json:"username"`
	Authorized bool   `json:"authorized"`
	jwt.StandardClaims
}
