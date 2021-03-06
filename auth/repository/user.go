package repository

import (
	"errors"
	"fmt"
	"go-todo/auth/data"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	addUserSchema    = `INSERT INTO users (username, password) VALUES (:username, :password)`
	authUserSchema   = `SELECT * FROM users WHERE username=?`
	deleteUserSchema = `DELETE FROM users where username=?`
)

func StoreUser(db *sqlx.DB, user data.User, logger *zap.SugaredLogger) (int64, error) {
	res, err := db.NamedExec(addUserSchema, user)
	if err != nil {
		logger.Error(err)
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		logger.Error(err)
		return -1, err
	}

	return id, nil
}

func AuthenticateUser(db *sqlx.DB, user data.User, logger *zap.SugaredLogger) (string, error) {
	row := db.QueryRow(authUserSchema, user.Username)
	if row.Err() != nil {
		logger.Error(row.Err().Error())
		return "", row.Err()
	}
	var username, password string
	var id int64
	err := row.Scan(&id, &username, &password)
	reqUser := data.User{Id: id, Username: username, Password: password}
	if err != nil {
		err = fmt.Errorf("invalid credentials")
		logger.Error(err)
		return "", err
	}

	isValid := user.ComparePassword(reqUser.Password)

	if !isValid {
		err = fmt.Errorf("invalid credentials")
		logger.Error(err)
		return "", err
	}

	token, err := reqUser.GetJWT()
	if err != nil {
		logger.Error(err)
		return "", err
	}

	return token, nil
}

func DeleteUser(db *sqlx.DB, username string, logger *zap.SugaredLogger) error {
	res, err := db.Exec(deleteUserSchema, username)
	if err != nil {
		logger.Error(err)
		return errors.New("could not delete user")
	}
	c, err := res.RowsAffected()
	if err != nil {
		logger.Error(err)
		return errors.New("could not delete user")
	}

	if c != 1 {
		m := fmt.Sprintf("Could not find user with username - %s", username)
		logger.Error(m)
		return errors.New(m)
	}

	return nil
}
