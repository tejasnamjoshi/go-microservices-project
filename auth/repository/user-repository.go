package repository

import "go-todo/auth/entities"

type UserRepository interface {
	Create(user *entities.User) (int64, error)
	Delete(username string) error
	Authenticate(user *entities.User) (*entities.User, error)
	GetAll() (*entities.Users, error)
}
