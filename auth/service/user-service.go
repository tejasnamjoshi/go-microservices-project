package service

import (
	"errors"
	"go-todo/auth/entities"
	"go-todo/auth/logging"
	"go-todo/auth/repository"

	"github.com/go-playground/validator/v10"
)

type UserService interface {
	Validate(user *entities.User) error
	Create(user *entities.User) error
	Login(user *entities.User) (string, error)
	Delete(username string) error
	GetAll() (*entities.Users, error)
}

type userServiceStruct struct {
	userRepository repository.UserRepository
	logger         logging.Logger
	jwtService     JWTService
}

func NewUserService(r repository.UserRepository, logger logging.Logger, jwtService JWTService) UserService {
	return &userServiceStruct{r, logger, jwtService}
}

func (*userServiceStruct) Validate(user *entities.User) error {
	v := validator.New()
	v.RegisterValidation("passwd", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > 5
	})
	return v.Struct(user)
}

func (us *userServiceStruct) Create(user *entities.User) error {
	err := us.jwtService.GeneratePassword(user)
	if err != nil {
		us.logger.Error(err.Error())
		return err
	}
	_, err = us.userRepository.Create(user)
	if err != nil {
		us.logger.Error(err.Error())
		return err
	}

	return nil
}

func (us *userServiceStruct) Login(user *entities.User) (string, error) {
	dbUser, err := us.userRepository.Authenticate(user)
	if err != nil {
		return "", err
	}
	isValid := us.jwtService.ComparePassword(user, dbUser)

	if !isValid {
		err = errors.New("invalid credentials")
		us.logger.Error(err.Error())
		return "", err
	}

	token, err := us.jwtService.GetJWT(dbUser)
	if err != nil {
		us.logger.Error(err.Error())
		return "", err
	}

	return token, nil
}

func (us *userServiceStruct) Delete(username string) error {
	return us.userRepository.Delete(username)
}

func (us *userServiceStruct) GetAll() (*entities.Users, error) {
	return us.userRepository.GetAll()
}
