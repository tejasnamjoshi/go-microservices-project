package repository

import (
	"database/sql"
	"errors"
	"go-todo/auth/entities"
	"go-todo/auth/logging"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func mockDb(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *sqlx.DB) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Failed to create mock db")
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	return mockDB, mock, sqlxDB
}

func Test_mysql_Create(t *testing.T) {
	_, mock, sqlxDB := mockDb(t)
	defer sqlxDB.Close()

	logger := zap.NewNop().Sugar()
	type fields struct {
		db     *sqlx.DB
		logger logging.Logger
	}
	type args struct {
		user *entities.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "success",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{user: &entities.User{Username: "test-username", Password: "test-password"}},
			want:    1,
			wantErr: false,
		},
		{
			name:    "error-exec",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{user: &entities.User{Username: "test-username", Password: "test-password"}},
			want:    0,
			wantErr: true,
		},
		{
			name:    "error-last-insertId",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{user: &entities.User{Username: "test-username", Password: "test-password"}},
			want:    0,
			wantErr: true,
		},
	}
	query := "INSERT INTO users (username, password) VALUES (?, ?)"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewMysqlRepository(tt.fields.db, tt.fields.logger)
			var expectedId int64 = 1

			if tt.name == "error-exec" {
				mock.ExpectExec(query).WithArgs(tt.args.user.Username, tt.args.user.Password).WillReturnError(errors.New("could not create user"))
			} else if tt.name == "error-last-insertId" {
				mock.ExpectExec(query).WithArgs(tt.args.user.Username, tt.args.user.Password).WillReturnResult(sqlmock.NewErrorResult(errors.New("could not fetch last inserted id")))
			} else {
				mock.ExpectExec(query).WithArgs(tt.args.user.Username, tt.args.user.Password).WillReturnResult(sqlmock.NewResult(expectedId, 1))
			}

			got, err := r.Create(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("mysql.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("mysql.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mysql_Delete(t *testing.T) {
	_, mock, sqlxDB := mockDb(t)
	defer sqlxDB.Close()

	logger := zap.NewNop().Sugar()
	type fields struct {
		db     *sqlx.DB
		logger logging.Logger
	}
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{username: "test-username"},
			wantErr: false,
		},
		{
			name:    "error-exec",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{username: "error-username"},
			wantErr: true,
		},
		{
			name:    "error-rows-affected",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{username: "error-username"},
			wantErr: true,
		},
		{
			name:    "error-rows-affected-0",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{username: "error-username"},
			wantErr: true,
		},
	}
	query := "DELETE FROM users where username=?"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewMysqlRepository(tt.fields.db,
				tt.fields.logger)

			if tt.name == "error-exec" {
				mock.ExpectExec(query).WithArgs(tt.args.username).WillReturnError(errors.New("could not delete user"))
			} else if tt.name == "error-rows-affected" {
				mock.ExpectExec(query).WithArgs(tt.args.username).WillReturnResult(sqlmock.NewErrorResult(errors.New("could not get rows affected")))
			} else if tt.name == "error-rows-affected-0" {
				mock.ExpectExec(query).WithArgs(tt.args.username).WillReturnResult(sqlmock.NewResult(1, 0))
			} else {
				mock.ExpectExec(query).WithArgs(tt.args.username).WillReturnResult(sqlmock.NewResult(1, 1))
			}

			if err := r.Delete(tt.args.username); (err != nil) != tt.wantErr {
				t.Errorf("mysql.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mysql_Authenticate(t *testing.T) {
	_, mock, sqlxDB := mockDb(t)
	defer sqlxDB.Close()

	logger := zap.NewNop().Sugar()

	type fields struct {
		db     *sqlx.DB
		logger logging.Logger
	}
	type args struct {
		user *entities.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entities.User
		wantErr bool
	}{
		{
			name:    "success",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{user: &entities.User{Username: "test-username", Password: "test-password"}},
			want:    &entities.User{Id: 1, Username: "test-username", Password: "test-password"},
			wantErr: false,
		},
		{
			name:    "error-row",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{user: &entities.User{Username: "test-username", Password: "test-password"}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error-scan",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{user: &entities.User{Username: "test-username", Password: "test-password"}},
			want:    nil,
			wantErr: true,
		},
	}
	query := "SELECT * FROM users WHERE username=?"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewMysqlRepository(tt.fields.db, tt.fields.logger)

			if tt.name == "error-scan" {
				mock.ExpectQuery(query).WithArgs(tt.args.user.Username).WillReturnError(errors.New("invalid credentials"))
			} else if tt.name == "error-row" {
				expectedRows := sqlmock.NewRows([]string{"id", "username", "password"})
				mock.ExpectQuery(query).WithArgs(tt.args.user.Username).WillReturnRows(expectedRows)
			} else {
				expectedRows := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow("1", tt.args.user.Username, tt.args.user.Password)
				mock.ExpectQuery(query).WithArgs(tt.args.user.Username).WillReturnRows(expectedRows)
			}

			got, err := r.Authenticate(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("mysql.Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mysql.Authenticate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mysql_GetAll(t *testing.T) {
	_, mock, sqlxDB := mockDb(t)
	defer sqlxDB.Close()

	logger := zap.NewNop().Sugar()

	type fields struct {
		db     *sqlx.DB
		logger logging.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		want    *entities.Users
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{db: sqlxDB, logger: logger},
			want: &entities.Users{
				entities.User{Id: 1, Username: "test-username", Password: "test-password"},
				entities.User{Id: 2, Username: "test-username-2", Password: "test-password-2"},
			},
			wantErr: false,
		},
		{
			name:    "error",
			fields:  fields{db: sqlxDB, logger: logger},
			want:    &entities.Users{},
			wantErr: true,
		},
	}
	query := "SELECT * FROM users"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewMysqlRepository(tt.fields.db, tt.fields.logger)

			if tt.name == "error" {
				mock.ExpectQuery(query).WillReturnError(errors.New("failed to fetch records"))
			} else {
				expectedRows := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow("1", "test-username", "test-password").AddRow("2", "test-username-2", "test-password-2")
				mock.ExpectQuery(query).WillReturnRows(expectedRows)
			}

			got, err := r.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("mysql.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mysql.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
