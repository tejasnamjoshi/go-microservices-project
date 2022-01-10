package repository

import (
	"errors"
	"fmt"
	"go-todo/auth/entities"
	"go-todo/auth/logging"

	"github.com/jmoiron/sqlx"
)

type mysql struct {
	db     *sqlx.DB
	logger logging.Logger
}

// Constructor function for the mysql repository.
func NewMysqlRepository(d *sqlx.DB, logger logging.Logger) UserRepository {
	return &mysql{d, logger}
}

// Inserts a user into the db and returns its id and error / nil
func (r *mysql) Create(user *entities.User) (int, error) {
	var addUserSchema = `INSERT INTO users (username, password) VALUES (:username, :password)`
	res, err := r.db.NamedExec(addUserSchema, user)
	if err != nil {
		r.logger.Error(err.Error())
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		r.logger.Error(err.Error())
		return 0, err
	}

	return int(id), nil
}

// Removes a user from the table and returns an error / nil
func (r *mysql) Delete(username string) error {
	var deleteUserSchema = `DELETE FROM users where username=?`
	res, err := r.db.Exec(deleteUserSchema, username)
	if err != nil {
		r.logger.Error(err.Error())
		return errors.New("could not delete user")
	}
	c, err := res.RowsAffected()
	if err != nil {
		r.logger.Error(err.Error())
		return errors.New("could not delete user")
	}

	if c != 1 {
		m := fmt.Sprintf("Could not find user with username - %s", username)
		r.logger.Error(m)
		return errors.New(m)
	}

	return nil
}

// Fetches a user record with the provided username and returns it
func (r *mysql) Authenticate(user *entities.User) (*entities.User, error) {
	var authUserSchema = `SELECT * FROM users WHERE username=?`
	row := r.db.QueryRow(authUserSchema, user.Username)
	if row.Err() != nil {
		r.logger.Error(row.Err().Error())
		return nil, row.Err()
	}
	var username, password string
	var id int
	err := row.Scan(&id, &username, &password)
	if err != nil {
		err = fmt.Errorf("invalid credentials")
		r.logger.Error(err.Error())
		return nil, err
	}
	dbUser := entities.User{Id: id, Username: username, Password: password}

	return &dbUser, nil
}

// Returns records of all users
func (r *mysql) GetAll() (*entities.Users, error) {
	var selectAllSchema = `SELECT * FROM users`
	users := entities.Users{}
	err := r.db.Select(&users, selectAllSchema)

	return &users, err
}
