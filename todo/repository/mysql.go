package repository

import (
	"errors"
	"go-todo/todo/entities"
	"go-todo/todo/logging"

	"github.com/jmoiron/sqlx"
)

type mysql struct {
	db     *sqlx.DB
	logger logging.Logger
}

func NewMysqlRepository(d *sqlx.DB, logger logging.Logger) TodoRepository {
	return &mysql{d, logger}
}
func (m *mysql) Create(todo *entities.Todo, userId int) error {
	createNewTodoSchema := `INSERT INTO todos (content) values (:content)`
	createNewUserTodoSchema := `INSERT INTO users_todos (user_id, todo_id) values (?, ?)`

	res, err := m.db.NamedExec(createNewTodoSchema, todo)
	if err != nil {
		m.logger.Error(err.Error())
		return err
	}

	todoId, err := res.LastInsertId()
	if err != nil {
		m.logger.Error(err.Error())
		return err
	}

	_, err = m.db.Exec(createNewUserTodoSchema, userId, todoId)
	if err != nil {
		m.logger.Error(err.Error())
		return err
	}
	return err
}

func (m *mysql) GetByUsername(userId int) (entities.Todos, error) {
	getTodoIdsByUserIdSchema := `SELECT todo_id from users_todos where user_id=?`
	getTodoByIds := "SELECT * from todos where id IN (?)"

	var todos = entities.Todos{}
	var todoIds = []int{}
	err := m.db.Select(&todoIds, getTodoIdsByUserIdSchema, userId)
	if err != nil {
		m.logger.Error(err.Error())
		return todos, err
	}

	if len(todoIds) > 0 {
		query, args, err := sqlx.In(getTodoByIds, todoIds)
		if err != nil {
			m.logger.Error(err.Error())
			return todos, err
		}
		query = m.db.Rebind(query)
		err = m.db.Select(&todos, query, args...)

		if err != nil {
			m.logger.Error(err.Error())
			return todos, err
		}
	}

	return todos, nil
}

func (m *mysql) MarkAsComplete(todoId string) error {
	markAsCompleteSchema := "UPDATE todos SET completed = 1 where id = ?"

	res, err := m.db.Exec(markAsCompleteSchema, todoId)
	if err != nil {
		m.logger.Error(err.Error())
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		m.logger.Error(err.Error())
		return err
	}

	if rowsAffected != 1 {
		err = errors.New("could not mark todo as completed")
		m.logger.Error(err.Error())
		return err
	}

	return nil
}
