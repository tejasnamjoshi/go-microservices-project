package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (t Todos) MarkAsComplete(rw http.ResponseWriter, r *http.Request) {
	todoId := chi.URLParam(r, "todoId")
	err := t.todoRepository.MarkAsComplete(todoId)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Could not complete the todo."))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Todo marked as completed."))
}
