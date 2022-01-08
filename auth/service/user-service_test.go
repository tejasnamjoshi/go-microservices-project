package service

import (
	"database/sql"
	"errors"
	"go-todo/auth/entities"
	"go-todo/auth/logging"
	"go-todo/auth/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

/* ********** MOCK ******************* */
type jwtservicemock struct {
	logger logging.Logger
}

func (*jwtservicemock) GeneratePassword(user *entities.User) error {
	if user.Password == "error-password" {
		return errors.New("mocked error")
	}
	return nil
}
func (*jwtservicemock) ComparePassword(user *entities.User, dbUser *entities.User) bool {
	return user.Password != "error-password"
}
func (*jwtservicemock) GetJWT(user *entities.User) (string, error) {
	if user.Username == "error-jwt" {
		return "error-token", errors.New("error generating token")
	}
	return "mock-token", nil
}
func (*jwtservicemock) GetAuthorizationData(authHeader string) (*CustomClaims, error) {
	return &CustomClaims{}, nil
}

/* ********** MOCK END ******************* */

func mockDb(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *sqlx.DB) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to create mock db")
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	return mockDB, mock, sqlxDB
}

func Test_userServiceStruct_Create(t *testing.T) {
	_, mock, sqlxDB := mockDb(t)
	logger := zap.NewNop().Sugar()
	userRepository := repository.NewMysqlRepository(sqlxDB, logger)
	jwtService := &jwtservicemock{logger: logger}

	type fields struct {
		userRepository repository.UserRepository
		logger         logging.Logger
		jwtService     JWTService
	}
	type args struct {
		user *entities.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{user: &entities.User{Username: "tejas", Password: "test-password"}},
			wantErr: false,
			fields: fields{
				logger:         logger,
				userRepository: userRepository,
				jwtService:     jwtService,
			},
		},
		{
			name:    "password error",
			args:    args{user: &entities.User{Username: "tejas", Password: "error-password"}},
			wantErr: true,
			fields: fields{
				logger:         logger,
				userRepository: userRepository,
				jwtService:     jwtService,
			},
		},
		{
			name:    "mysql-error",
			args:    args{user: &entities.User{Username: "tejas", Password: "test-password"}},
			wantErr: true,
			fields: fields{
				logger:         logger,
				userRepository: userRepository,
				jwtService:     jwtService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := NewUserService(
				tt.fields.userRepository,
				tt.fields.logger,
				tt.fields.jwtService,
			)
			if tt.name == "mysql-error" {
				mock.ExpectExec(`INSERT`).WithArgs(tt.args.user.Username, tt.args.user.Password).WillReturnError(errors.New("error while inserting record"))
			} else {
				mock.ExpectExec(`INSERT`).WithArgs(tt.args.user.Username, tt.args.user.Password).WillReturnResult(sqlmock.NewResult(int64(1), 1))
			}

			if err := us.Create(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userServiceStruct.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userServiceStruct_Login(t *testing.T) {
	_, mock, sqlxDB := mockDb(t)
	logger := zap.NewNop().Sugar()
	userRepository := repository.NewMysqlRepository(sqlxDB, logger)
	jwtService := &jwtservicemock{logger: logger}

	type fields struct {
		userRepository repository.UserRepository
		logger         logging.Logger
		jwtService     JWTService
	}
	type args struct {
		user *entities.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{user: &entities.User{Username: "tejas", Password: "test-password"}},
			wantErr: false,
			fields: fields{
				logger:         logger,
				userRepository: userRepository,
				jwtService:     jwtService,
			},
			want: "mock-token",
		},
		{
			name:    "error-authenticate",
			args:    args{user: &entities.User{Username: "tejas", Password: "error-password"}},
			wantErr: true,
			fields: fields{
				logger:         logger,
				userRepository: userRepository,
				jwtService:     jwtService,
			},
			want: "",
		},
		{
			name:    "error-password",
			args:    args{user: &entities.User{Username: "tejas", Password: "error-password"}},
			wantErr: true,
			fields: fields{
				logger:         logger,
				userRepository: userRepository,
				jwtService:     jwtService,
			},
			want: "",
		},
		{
			name:    "error-jwt",
			args:    args{user: &entities.User{Username: "error-jwt", Password: "test-password"}},
			wantErr: true,
			fields: fields{
				logger:         logger,
				userRepository: userRepository,
				jwtService:     jwtService,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := NewUserService(
				tt.fields.userRepository,
				tt.fields.logger,
				tt.fields.jwtService,
			)

			if tt.name == "error-authenticate" {
				mock.ExpectQuery("SELECT").WithArgs(tt.args.user.Username).WillReturnError(errors.New("failed to authenticate"))
			} else {
				expectedRows := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow("1", tt.args.user.Username, tt.args.user.Password)
				mock.ExpectQuery("SELECT").WithArgs(tt.args.user.Username).WillReturnRows(expectedRows)
			}

			got, err := us.Login(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("userServiceStruct.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("userServiceStruct.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userServiceStruct_Validate(t *testing.T) {
	_, _, sqlxDB := mockDb(t)
	logger := zap.NewNop().Sugar()
	userRepository := repository.NewMysqlRepository(sqlxDB, logger)
	jwtService := &jwtservicemock{logger: logger}

	type fields struct {
		userRepository repository.UserRepository
		logger         logging.Logger
		jwtService     JWTService
	}
	type args struct {
		user *entities.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				logger:         logger,
				userRepository: userRepository,
				jwtService:     jwtService,
			},
			args:    args{user: &entities.User{Username: "test-username", Password: "test-password"}},
			wantErr: false,
		},
		{
			name: "error-blank-username",
			fields: fields{
				logger:         logger,
				userRepository: userRepository,
				jwtService:     jwtService,
			},
			args:    args{user: &entities.User{Username: "", Password: "test-password"}},
			wantErr: true,
		},
		{
			name: "error-blank-password",
			fields: fields{
				logger:         logger,
				userRepository: userRepository,
				jwtService:     jwtService,
			},
			args:    args{user: &entities.User{Username: "test-username", Password: ""}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := NewUserService(
				tt.fields.userRepository,
				tt.fields.logger,
				tt.fields.jwtService,
			)
			if err := us.Validate(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userServiceStruct.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
