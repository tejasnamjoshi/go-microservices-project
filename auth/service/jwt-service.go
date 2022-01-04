package service

import (
	"errors"
	"go-todo/auth/entities"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type JWTService interface {
	GeneratePassword(user *entities.User) error
	ComparePassword(user *entities.User, dbUser *entities.User) bool
	GetJWT(user *entities.User) (string, error)
	GetAuthorizationData(authHeader string) (*CustomClaims, error)
}

type JWTServiceStruct struct{}

type CustomClaims struct {
	UserId     int64  `json:"userId"`
	Username   string `json:"username"`
	Authorized bool   `json:"authorized"`
	jwt.StandardClaims
}

func NewJWTService() JWTService {
	return &JWTServiceStruct{}
}

func (*JWTServiceStruct) GeneratePassword(user *entities.User) error {
	pwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	user.Password = string(pwd)
	return nil
}

func (*JWTServiceStruct) ComparePassword(user *entities.User, dbUser *entities.User) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	return err == nil
}

func (*JWTServiceStruct) GetJWT(user *entities.User) (string, error) {
	var mySigningKey = []byte(os.Getenv("SECRET_KEY"))
	st, _ := strconv.Atoi(os.Getenv("SESSION_TIME"))
	claims := CustomClaims{
		UserId:     user.Id,
		Username:   user.Username,
		Authorized: true,
		StandardClaims: jwt.StandardClaims{
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
	var mySigningKey = []byte(os.Getenv("SECRET_KEY"))

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, errors.New(("invalid Signing Method"))
		}
		aud := os.Getenv("JWT_AUD")
		checkAudience := token.Claims.(*CustomClaims).VerifyAudience(aud, false)
		if !checkAudience {
			return false, errors.New(("invalid aud"))
		}
		iss := os.Getenv("JWT_ISS")
		checkIss := token.Claims.(*CustomClaims).VerifyIssuer(iss, false)
		if !checkIss {
			return false, errors.New(("invalid iss"))
		}
		return mySigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("error decoding token")
	}

	return claims, nil
}

func (*JWTServiceStruct) GetAuthorizationData(userToken string) (*CustomClaims, error) {
	if userToken == "" {
		return nil, errors.New("no Authorization Token provided")
	}

	return ParseJWT(userToken)
}
