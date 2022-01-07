package service

import (
	"errors"
	"go-todo/auth/entities"
	"go-todo/auth/logging"
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

type userServiceStruct struct {
	logger logging.Logger
}

func NewUserService(r repository.UserRepository, logger logging.Logger) UserService {
	userRepository = r
	return &userServiceStruct{logger}
}

func (*userServiceStruct) Validate(user *entities.User) error {
	v := validator.New()
	v.RegisterValidation("passwd", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > 5
	})
	return v.Struct(user)
}

func (us *userServiceStruct) Create(user *entities.User) error {
	jwtService := NewJWTService(us.logger)
	err := jwtService.GeneratePassword(user)
	if err != nil {
		us.logger.Error(err.Error())
		return err
	}
	_, err = userRepository.Create(user)
	if err != nil {
		us.logger.Error(err.Error())
		return err
	}

	return nil
}

func (us *userServiceStruct) Login(user *entities.User) (string, error) {
	dbUser, err := userRepository.Authenticate(user)
	if err != nil {
		return "", err
	}
	jwtService := NewJWTService(us.logger)
	isValid := jwtService.ComparePassword(user, dbUser)

	if !isValid {
		err = errors.New("invalid credentials")
		us.logger.Error(err.Error())
		return "", err
	}

	token, err := jwtService.GetJWT(dbUser)
	if err != nil {
		us.logger.Error(err.Error())
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
