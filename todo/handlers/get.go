package handlers

import (
	"net/http"

	"go-todo/todo/data"
	todo_db "go-todo/todo/db"

	"github.com/jmoiron/sqlx"
)

var getTodoIdsByUserIdSchema = `Select todo_id from users_todos where user_id=?`
var getTodoByIds = "Select * from todos where id IN (?)"

func (t Todos) GetByUsername(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := ctx.Value("userId").(int64)
	var todos = data.Todos{}
	var todoIds = []int{}
	var db = todo_db.GetDb()
	err := db.Select(&todoIds, getTodoIdsByUserIdSchema, id)
	if err != nil {
		t.l.Panicln(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(todoIds) > 0 {
		query, args, err := sqlx.In(getTodoByIds, todoIds)
		if err != nil {
			t.l.Panicln(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		query = db.Rebind(query)
		err = db.Select(&todos, query, args...)

		if err != nil {
			t.l.Panicln(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	todos.ToJSON(rw)
}
