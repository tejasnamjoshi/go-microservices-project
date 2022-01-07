package service

import (
	"go-todo/todo/entities"
	"go-todo/todo/logging"
	"go-todo/todo/repository"
)

type TodoService interface {
	Create(todo *entities.Todo, userId int) error
	GetByUsername(userId int) (entities.Todos, error)
	MarkAsComplete(todoId string) error
}

type todoservice struct {
	todoRepository repository.TodoRepository
	logger         logging.Logger
}

func NewTodoService(todoRepository repository.TodoRepository, logger logging.Logger) TodoService {
	return &todoservice{todoRepository, logger}
}

func (t *todoservice) Create(todo *entities.Todo, userId int) error {
	return t.todoRepository.Create(todo, userId)
}
func (t *todoservice) GetByUsername(userId int) (entities.Todos, error) {
	return t.todoRepository.GetByUsername(userId)
}
func (t *todoservice) MarkAsComplete(todoId string) error {
	return t.todoRepository.MarkAsComplete(todoId)
}
