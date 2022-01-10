package service

import (
	"errors"
	"go-todo/auth/entities"
	"go-todo/auth/logging"
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

type JWTServiceStruct struct {
	logger logging.Logger
}

type CustomClaims struct {
	UserId     int    `json:"userId"`
	Username   string `json:"username"`
	Authorized bool   `json:"authorized"`
	jwt.StandardClaims
}

func NewJWTService(logger logging.Logger) JWTService {
	return &JWTServiceStruct{logger}
}

// Generates a password using the brypt hashing function and adds it to the user entity
// Returns an error or nil
func (js *JWTServiceStruct) GeneratePassword(user *entities.User) error {
	pwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		js.logger.Error(err.Error())
		return err
	}

	user.Password = string(pwd)
	return nil
}

// Returns true if hash and plain-text password match, else returns false
func (*JWTServiceStruct) ComparePassword(user *entities.User, dbUser *entities.User) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	return err == nil
}

// Generates a JWT using custom claims as defined in the .env file.
// Returns token and error / nil
func (js *JWTServiceStruct) GetJWT(user *entities.User) (string, error) {
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
		js.logger.Error(err.Error())
		return "", err
	}

	return tokenString, nil
}

// Parses the provider JWT and returns claims and error / nil
func ParseJWT(tokenString string, logger logging.Logger) (*CustomClaims, error) {
	var mySigningKey = []byte(os.Getenv("SECRET_KEY"))

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := errors.New(("invalid Signing Method"))
			logger.Error(err.Error())
			return false, err
		}
		aud := os.Getenv("JWT_AUD")
		checkAudience := token.Claims.(*CustomClaims).VerifyAudience(aud, false)
		if !checkAudience {
			err := errors.New(("invalid aud"))
			logger.Error(err.Error())
			return false, err
		}
		iss := os.Getenv("JWT_ISS")
		checkIss := token.Claims.(*CustomClaims).VerifyIssuer(iss, false)
		if !checkIss {
			err := errors.New(("invalid iss"))
			logger.Error(err.Error())
			return false, err
		}
		return mySigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		err := errors.New(("error decoding token"))
		logger.Error(err.Error())
		return nil, err
	}

	return claims, nil
}

// Returns the response received from ParseJWT function above
func (js *JWTServiceStruct) GetAuthorizationData(userToken string) (*CustomClaims, error) {
	if userToken == "" {
		err := errors.New("no Authorization Token provided")
		js.logger.Error(err.Error())
		return nil, err
	}

	return ParseJWT(userToken, js.logger)
}
