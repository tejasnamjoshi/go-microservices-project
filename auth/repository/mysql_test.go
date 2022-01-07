package repository_test

import (
	"database/sql"
	"errors"
	"go-todo/auth/entities"
	"go-todo/auth/logging"
	"go-todo/auth/repository"
	"go-todo/auth/service"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

var (
	userRepository repository.UserRepository
	logger         logging.Logger     = logging.NewZapLogger()
	jwtService     service.JWTService = service.NewJWTService(logger)
)

func mockDb(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to create mock db")
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	userRepository = repository.NewMysqlRepository(sqlxDB, logger)
	return mockDB, mock
}

func TestStoreUser(t *testing.T) {
	mockDB, mock := mockDb(t)
	defer mockDB.Close()
	var expectedId int = 1
	user := entities.User{Username: "test-user", Password: "test-password"}

	mock.ExpectExec(`INSERT INTO users`).WithArgs(user.Username, user.Password).WillReturnResult(sqlmock.NewResult(int64(expectedId), 1))

	id, err := userRepository.Create(&user)
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
	mockDB, mock := mockDb(t)
	defer mockDB.Close()
	var expectedId int = 1
	user := entities.User{Username: "test-user", Password: "test-password"}

	mock.ExpectExec(`INSERT INTO users`).WithArgs(user.Username, user.Password).WillReturnError(errors.New("failed to insert"))

	id, err := userRepository.Create(&user)
	if err == nil || id != -1 {
		t.Errorf("Expected user with id %d but got error %s", expectedId, err)
	}
}

func TestLogin(t *testing.T) {
	mockDB, mock := mockDb(t)
	defer mockDB.Close()
	var username, password string = "test-user", "test-password"

	user := entities.User{Username: username, Password: password}
	err := jwtService.GeneratePassword(&user)
	if err != nil {
		t.Errorf("Could not generate password : %v", err)
	}

	expectedRows := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow("1", username, password)
	mock.ExpectQuery("SELECT").WithArgs(username).WillReturnRows(expectedRows)

	_, err = userRepository.Authenticate(&user)
	if err != nil {
		t.Errorf("Expected token but got %v", err)
	}
}

func TestLoginError(t *testing.T) {
	mockDB, mock := mockDb(t)
	defer mockDB.Close()
	var username, password string = "test-user", "test-password"

	user := entities.User{Username: username, Password: password}
	err := jwtService.GeneratePassword(&user)
	if err != nil {
		t.Errorf("Could not generate password : %v", err)
	}
	mock.ExpectQuery("SELECT").WithArgs(username).WillReturnError(errors.New("invalid credentials"))

	_, err = userRepository.Authenticate(&user)
	if err == nil {
		t.Errorf("Expected token but got %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	mockDB, mock := mockDb(t)
	defer mockDB.Close()
	var username string = "test-user"

	mock.ExpectExec("DELETE").WithArgs(username).WillReturnResult(sqlmock.NewResult(1, 1))
	err := userRepository.Delete(username)
	if err != nil {
		t.Errorf("Expected nil error but got %v", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not fulfilled: %s", err)
	}
}

func TestDeleteUserError(t *testing.T) {
	mockDB, mock := mockDb(t)
	defer mockDB.Close()
	var username string = "test-user"

	mock.ExpectExec("DELETE").WithArgs(username).WillReturnError(errors.New("could not delete user"))
	err := userRepository.Delete(username)
	if err == nil {
		t.Errorf("Expected nil error but got %v", err)
	}
}
