package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (t Todos) MarkAsComplete(rw http.ResponseWriter, r *http.Request) {
	todoId := chi.URLParam(r, "todoId")
	err := t.TodoService.MarkAsComplete(todoId)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Could not complete the todo."))
		return
	}

	t.Response.CreateSuccessResponse(rw, "Todo marked as completed")
}
