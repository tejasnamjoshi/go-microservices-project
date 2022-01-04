package repository_test

import (
	"database/sql"
	"errors"
	"go-todo/auth/data"
	"go-todo/auth/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func mockDb(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *sqlx.DB) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to create mock db")
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	return mockDB, mock, sqlxDB
}

func TestStoreUser(t *testing.T) {
	mockDB, mock, sqlxDB := mockDb(t)
	defer mockDB.Close()
	var expectedId int64 = 1
	user := data.User{Username: "test-user", Password: "test-password"}

	mock.ExpectExec(`INSERT INTO users`).WithArgs(user.Username, user.Password).WillReturnResult(sqlmock.NewResult(expectedId, 1))

	logger := zap.NewNop().Sugar()
	id, err := repository.StoreUser(sqlxDB, user, logger)
	if err != nil {
		t.Errorf("Expected user with id %d but got error %s", expectedId, err)
	}

	if id != expectedId {
		t.Errorf("Expected user with id %d but got %d", expectedId, id)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not fulfilled: %s", err)
	}
}

func TestStoreUserError(t *testing.T) {
	mockDB, mock, sqlxDB := mockDb(t)
	defer mockDB.Close()
	var expectedId int64 = 1
	user := data.User{Username: "test-user", Password: "test-password"}

	mock.ExpectExec(`INSERT INTO users`).WithArgs(user.Username, user.Password).WillReturnError(errors.New("failed to insert"))

	logger := zap.NewNop().Sugar()
	id, err := repository.StoreUser(sqlxDB, user, logger)
	if err == nil || id != -1 {
		t.Errorf("Expected user with id %d but got error %s", expectedId, err)
	}
}

func TestLogin(t *testing.T) {
	mockDB, mock, sqlxDB := mockDb(t)
	defer mockDB.Close()
	var username, password string = "test-user", "test-password"

	user := data.User{Username: username, Password: password}
	err := user.GeneratePassword()
	if err != nil {
		t.Errorf("Could not generate password : %v", err)
	}
	logger := zap.NewNop().Sugar()

	expectedRows := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow("1", username, password)
	mock.ExpectQuery("SELECT").WithArgs(username).WillReturnRows(expectedRows)

	_, err = repository.AuthenticateUser(sqlxDB, user, logger)
	if err != nil {
		t.Errorf("Expected token but got %v", err)
	}
}

func TestLoginError(t *testing.T) {
	mockDB, mock, sqlxDB := mockDb(t)
	defer mockDB.Close()
	var username, password string = "test-user", "test-password"

	user := data.User{Username: username, Password: password}
	err := user.GeneratePassword()
	if err != nil {
		t.Errorf("Could not generate password : %v", err)
	}
	logger := zap.NewNop().Sugar()

	// expectedRows := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow("1", username, password)
	mock.ExpectQuery("SELECT").WithArgs(username).WillReturnError(errors.New("invalid credentials"))

	_, err = repository.AuthenticateUser(sqlxDB, user, logger)
	if err == nil {
		t.Errorf("Expected token but got %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	mockDB, mock, sqlxDB := mockDb(t)
	defer mockDB.Close()
	logger := zap.NewNop().Sugar()
	var username string = "test-user"

	mock.ExpectExec("DELETE").WithArgs(username).WillReturnResult(sqlmock.NewResult(1, 1))
	err := repository.DeleteUser(sqlxDB, username, logger)
	if err != nil {
		t.Errorf("Expected nil error but got %v", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not fulfilled: %s", err)
	}
}

func TestDeleteUserError(t *testing.T) {
	mockDB, mock, sqlxDB := mockDb(t)
	defer mockDB.Close()
	logger := zap.NewNop().Sugar()
	var username string = "test-user"

	mock.ExpectExec("DELETE").WithArgs(username).WillReturnError(errors.New("could not delete user"))
	err := repository.DeleteUser(sqlxDB, username, logger)
	if err == nil {
		t.Errorf("Expected nil error but got %v", err)
	}
}
