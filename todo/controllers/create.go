package controllers

import (
	"go-todo/todo/entities"
	"net/http"
)

func (t Todos) CreateNewTodo(rw http.ResponseWriter, r *http.Request) {
	var todo = entities.Todo{}
	ctx := r.Context()
	userId := ctx.Value(UserIdContext{}).(int)

	defer r.Body.Close()
	err := todo.FromJSON(r.Body)
	if err != nil {
		t.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = t.TodoService.Create(&todo, userId)
	if err != nil {
		t.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("TODO added successfully"))
}
