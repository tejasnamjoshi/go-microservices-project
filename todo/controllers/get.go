package controllers

import (
	"net/http"
)

func (t Todos) GetByUsername(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value(UserIdContext{}).(int)
	todos, err := t.TodoService.GetByUsername(userId)

	if err != nil {
		t.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	t.Response.CreateSuccessResponse(rw, todos)
}
