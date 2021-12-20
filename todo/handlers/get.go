package handlers

import (
	"fmt"
	"net/http"

	"go-todo/todo/data"
	todo_db "go-todo/todo/db"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

var getTodoIdsByUserIdSchema = `Select todo_id from users_todos where user_id=?`
var getTodoByIds = "Select * from todos where id IN (?)"

func (t Todos) GetByUsername(rw http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "username")
	var todos = data.Todos{}
	var todoIds = []int{}
	var db = todo_db.GetDb()

	err := db.Select(&todoIds, getTodoIdsByUserIdSchema, id)
	if err != nil {
		t.l.Panicln(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	query, args, err := sqlx.In(getTodoByIds, todoIds)
	if err != nil {
		t.l.Panicln(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	query = db.Rebind(query)
	fmt.Println(query, args)
	err = db.Select(&todos, query, args...)

	if err != nil {
		t.l.Panicln(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	todos.ToJSON(rw)
}
