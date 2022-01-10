package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (t Todos) MarkAsComplete(rw http.ResponseWriter, r *http.Request) {
	// Extract and Format the request data
	todoId := chi.URLParam(r, "todoId")

	// Invoke Logic
	err := t.TodoService.MarkAsComplete(todoId)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Could not complete the todo."))
		return
	}

	// Send Response
	t.Response.SendSuccessResponse(rw, "Todo marked as completed")
}
