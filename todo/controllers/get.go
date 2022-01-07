package controllers

import (
	"net/http"
)

func (t Todos) GetByUsername(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value("userId").(int)
	todos, err := t.TodoService.GetByUsername(userId)

	if err != nil {
		t.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	todos.ToJSON(rw)
}
