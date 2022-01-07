package repository

import (
	"errors"
	"fmt"
	"go-todo/auth/entities"

	"github.com/jmoiron/sqlx"
)

const (
	addUserSchema    = `INSERT INTO users (username, password) VALUES (:username, :password)`
	authUserSchema   = `SELECT * FROM users WHERE username=?`
	deleteUserSchema = `DELETE FROM users where username=?`
	selectAllSchema  = `SELECT * FROM users`
)

type repo struct{}

var db *sqlx.DB

func NewMysqlRepository(d *sqlx.DB) UserRepository {
	db = d
	return &repo{}
}

func (r *repo) Create(user *entities.User) (int, error) {
	res, err := db.NamedExec(addUserSchema, user)
	if err != nil {
		// logger.Error(err)
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		// logger.Error(err)
		return 0, err
	}

	return int(id), nil
}

func (r *repo) Delete(username string) error {
	res, err := db.Exec(deleteUserSchema, username)
	if err != nil {
		// logger.Error(err)
		return errors.New("could not delete user")
	}
	c, err := res.RowsAffected()
	if err != nil {
		// logger.Error(err)
		return errors.New("could not delete user")
	}

	if c != 1 {
		m := fmt.Sprintf("Could not find user with username - %s", username)
		// logger.Error(m)
		return errors.New(m)
	}

	return nil
}

func (r *repo) Authenticate(user *entities.User) (*entities.User, error) {
	row := db.QueryRow(authUserSchema, user.Username)
	if row.Err() != nil {
		// logger.Error(row.Err().Error())
		return nil, row.Err()
	}
	var username, password string
	var id int
	err := row.Scan(&id, &username, &password)
	dbUser := entities.User{Id: id, Username: username, Password: password}
	if err != nil {
		err = fmt.Errorf("invalid credentials")
		// logger.Error(err)
		return nil, err
	}

	return &dbUser, nil
}

func (r *repo) GetAll() (*entities.Users, error) {
	users := entities.Users{}
	err := db.Select(&users, selectAllSchema)

	return &users, err
}
