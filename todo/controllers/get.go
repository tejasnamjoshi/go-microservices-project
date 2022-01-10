package controllers

import (
	"net/http"
)

func (t Todos) GetByUsername(rw http.ResponseWriter, r *http.Request) {
	// Extract and Format the request data
	ctx := r.Context()
	userId := ctx.Value(UserIdContext{}).(int)

	// Invoke logic
	todos, err := t.TodoService.GetByUsername(userId)

	if err != nil {
		t.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send Response
	t.Response.SendSuccessResponse(rw, todos)
}
