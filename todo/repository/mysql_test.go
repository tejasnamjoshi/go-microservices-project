package repository

import (
	"database/sql"
	"errors"
	"go-todo/todo/entities"
	"go-todo/todo/logging"
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
		todo   *entities.Todo
		userId int
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
			args:    args{todo: &entities.Todo{Content: "This is a test todo", Completed: false}, userId: 1},
			wantErr: false,
		},
		{
			name:    "error-insert-todo",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{todo: &entities.Todo{Content: "This is a test todo", Completed: false}, userId: 1},
			wantErr: true,
		},
		{
			name:    "error-insert-todo-last-insert-id",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{todo: &entities.Todo{Content: "This is a test todo", Completed: false}, userId: 1},
			wantErr: true,
		},
		{
			name:    "error-insert-user-todo",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{todo: &entities.Todo{Content: "This is a test todo", Completed: false}, userId: 1},
			wantErr: true,
		},
	}
	query1 := "INSERT INTO todos (content) values (?)"
	query2 := "INSERT INTO users_todos (user_id, todo_id) values (?, ?)"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMysqlRepository(tt.fields.db, tt.fields.logger)
			var insertedTodoId int64 = 1

			if tt.name == "error-insert-todo" {
				mock.ExpectExec(query1).WithArgs(tt.args.todo.Content).WillReturnError(errors.New("failed to insert record"))
			} else if tt.name == "error-insert-todo-last-insert-id" {
				mock.ExpectExec(query1).WithArgs(tt.args.todo.Content).WillReturnResult(sqlmock.NewErrorResult(errors.New("could not get last inserted id")))
			} else if tt.name == "error-insert-user-todo" {
				mock.ExpectExec(query1).WithArgs(tt.args.todo.Content).WillReturnResult(sqlmock.NewResult(insertedTodoId, 1))
				mock.ExpectExec(query2).WithArgs(tt.args.userId, insertedTodoId).WillReturnError(errors.New("could not insert record"))
			} else {
				mock.ExpectExec(query1).WithArgs(tt.args.todo.Content).WillReturnResult(sqlmock.NewResult(insertedTodoId, 1))
				mock.ExpectExec(query2).WithArgs(tt.args.userId, insertedTodoId).WillReturnResult(sqlmock.NewResult(1, 1))
			}

			if err := m.Create(tt.args.todo, tt.args.userId); (err != nil) != tt.wantErr {
				t.Errorf("mysql.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mysql_GetByUsername(t *testing.T) {
	_, mock, sqlxDB := mockDb(t)
	defer sqlxDB.Close()

	logger := zap.NewNop().Sugar()

	type fields struct {
		db     *sqlx.DB
		logger logging.Logger
	}
	type args struct {
		userId int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entities.Todos
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{db: sqlxDB, logger: logger},
			args:   args{userId: 1},
			want: entities.Todos{
				entities.Todo{Id: 1, Content: "This is the first test todo", Completed: false},
				entities.Todo{Id: 2, Content: "This is the second test todo", Completed: true},
			},
			wantErr: false,
		},
		{
			name:    "success-0-todos",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{userId: 1},
			want:    entities.Todos{},
			wantErr: false,
		},
		{
			name:    "error-todoid",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{userId: 1},
			want:    entities.Todos{},
			wantErr: true,
		},
		{
			name:    "error-todo",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{userId: 1},
			want:    entities.Todos{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		expectedTodos := entities.Todos{
			entities.Todo{Id: 1, Content: "This is the first test todo", Completed: false},
			entities.Todo{Id: 2, Content: "This is the second test todo", Completed: true},
		}
		query1 := "SELECT todo_id from users_todos where user_id=?"
		query2 := "SELECT * from todos where id IN (?, ?)"

		t.Run(tt.name, func(t *testing.T) {
			m := NewMysqlRepository(tt.fields.db, tt.fields.logger)

			if tt.name == "error-todoid" {
				mock.ExpectQuery(query1).WithArgs(tt.args.userId).WillReturnError(errors.New("failed to get records"))
			} else if tt.name == "success-0-todos" {
				todoIds_rows := sqlmock.NewRows([]string{"todo_id"})
				mock.ExpectQuery(query1).WithArgs(tt.args.userId).WillReturnRows(todoIds_rows)
			} else if tt.name == "error-todo" {
				todoIds_rows := sqlmock.NewRows([]string{"todo_id"}).AddRow(expectedTodos[0].Id).AddRow(expectedTodos[1].Id)
				mock.ExpectQuery(query1).WithArgs(tt.args.userId).WillReturnRows(todoIds_rows)

				mock.ExpectQuery(query2).WithArgs(expectedTodos[0].Id, expectedTodos[1].Id).WillReturnError(errors.New("failed to fetch records"))
			} else {
				todoIds_rows := sqlmock.NewRows([]string{"todo_id"}).AddRow(expectedTodos[0].Id).AddRow(expectedTodos[1].Id)
				mock.ExpectQuery(query1).WithArgs(tt.args.userId).WillReturnRows(todoIds_rows)

				todos_rows := sqlmock.NewRows([]string{"id", "content", "completed"}).AddRow(expectedTodos[0].Id, expectedTodos[0].Content, expectedTodos[0].Completed).AddRow(expectedTodos[1].Id, expectedTodos[1].Content, expectedTodos[1].Completed)
				mock.ExpectQuery(query2).WithArgs(expectedTodos[0].Id, expectedTodos[1].Id).WillReturnRows(todos_rows)
			}

			got, err := m.GetByUsername(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("mysql.GetByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mysql.GetByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mysql_MarkAsComplete(t *testing.T) {
	_, mock, sqlxDB := mockDb(t)
	defer sqlxDB.Close()

	logger := zap.NewNop().Sugar()

	type fields struct {
		db     *sqlx.DB
		logger logging.Logger
	}
	type args struct {
		todoId string
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
			args:    args{todoId: "1"},
			wantErr: false,
		},
		{
			name:    "error",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{todoId: "1"},
			wantErr: true,
		},
		{
			name:    "error-rows-affected",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{todoId: "1"},
			wantErr: true,
		},
		{
			name:    "error-rows-affected-0",
			fields:  fields{db: sqlxDB, logger: logger},
			args:    args{todoId: "1"},
			wantErr: true,
		},
	}
	query := "UPDATE todos SET completed = 1 where id = ?"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMysqlRepository(tt.fields.db, tt.fields.logger)

			if tt.name == "error" {
				mock.ExpectExec(query).WithArgs(tt.args.todoId).WillReturnError(errors.New("failed to update record"))
			} else if tt.name == "error-rows-affected-0" {
				mock.ExpectExec(query).WithArgs(tt.args.todoId).WillReturnResult(sqlmock.NewResult(1, 0))
			} else if tt.name == "error-rows-affected" {
				mock.ExpectExec(query).WithArgs(tt.args.todoId).WillReturnResult(sqlmock.NewErrorResult(errors.New("failed to fetch rows affected")))
			} else {
				mock.ExpectExec(query).WithArgs(tt.args.todoId).WillReturnResult(sqlmock.NewResult(1, 1))
			}

			if err := m.MarkAsComplete(tt.args.todoId); (err != nil) != tt.wantErr {
				t.Errorf("mysql.MarkAsComplete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
