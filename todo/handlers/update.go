package handlers

import (
	todo_db "go-todo/todo/db"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var markAsCompleteSchema = "UPDATE todos SET completed = 1 where id = ?"

func (t Todos) MarkAsComplete(rw http.ResponseWriter, r *http.Request) {
	todoId := chi.URLParam(r, "todoId")
	res, err := todo_db.GetDb().Exec(markAsCompleteSchema, todoId)
	if err != nil {
		t.l.Panicln(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		t.l.Panicln(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rowsAffected != 1 {
		t.l.Panicln(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Could not complete the todo."))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Todo marked as completed."))
}
