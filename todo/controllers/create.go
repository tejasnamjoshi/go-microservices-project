package controllers

import (
	"go-todo/todo/entities"
	"net/http"
)

func (t Todos) CreateNewTodo(rw http.ResponseWriter, r *http.Request) {
	// Extract and Format the request data
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

	// Invoke logic
	err = t.TodoService.Create(&todo, userId)
	if err != nil {
		t.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send Response
	t.Response.SendSuccessResponse(rw, "TODO added successfully")
}
