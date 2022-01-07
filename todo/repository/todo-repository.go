package repository

import "go-todo/todo/entities"

type TodoRepository interface {
	Create(todo *entities.Todo, userId int) error
	GetByUsername(userId int) (entities.Todos, error)
	MarkAsComplete(todoId string) error
}
