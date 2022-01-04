package service

import (
	"errors"
	"go-todo/auth/entities"
	"go-todo/auth/repository"

	"github.com/go-playground/validator/v10"
)

var (
	userRepository repository.UserRepository
)

type UserService interface {
	Validate(user *entities.User) error
	Create(user *entities.User) error
	Login(user *entities.User) (string, error)
	Delete(username string) error
	GetAll() (*entities.Users, error)
}

type userServiceStruct struct{}

func NewUserService(r repository.UserRepository) UserService {
	userRepository = r
	return &userServiceStruct{}
}

func (*userServiceStruct) Validate(user *entities.User) error {
	v := validator.New()
	v.RegisterValidation("passwd", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > 5
	})
	return v.Struct(user)
}

func (*userServiceStruct) Create(user *entities.User) error {
	jwtService := NewJWTService()
	err := jwtService.GeneratePassword(user)
	if err != nil {
		return err
	}
	_, err = userRepository.Create(user)
	if err != nil {
		return err
	}

	return nil
}

func (*userServiceStruct) Login(user *entities.User) (string, error) {
	dbUser, err := userRepository.Authenticate(user)
	if err != nil {
		return "", err
	}
	jwtService := NewJWTService()
	isValid := jwtService.ComparePassword(user, dbUser)

	if !isValid {
		err = errors.New("invalid credentials")
		// logger.Error(err)
		return "", err
	}

	token, err := jwtService.GetJWT(user)
	if err != nil {
		// logger.Error(err)
		return "", err
	}

	return token, nil
}

func (*userServiceStruct) Delete(username string) error {
	return userRepository.Delete(username)
}

func (*userServiceStruct) GetAll() (*entities.Users, error) {
	users, err := userRepository.GetAll()

	return users, err
}
