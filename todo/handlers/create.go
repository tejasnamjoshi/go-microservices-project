package handlers

import (
	"go-todo/todo/data"
	todo_db "go-todo/todo/db"
	"net/http"
)

var createNewTodoSchema = `INSERT INTO todos (content) values (:content)`
var createNewUserTodoSchema = `INSERT INTO users_todos (user_id, todo_id) values (?, ?)`

func (t Todos) CreateNewTodo(rw http.ResponseWriter, r *http.Request) {
	var todo = data.Todo{}
	ctx := r.Context()
	userId := ctx.Value("userId").(int64)
	db := todo_db.GetDb()

	defer r.Body.Close()
	err := todo.FromJSON(r.Body)
	if err != nil {
		HandleError(err, rw, t)
		return
	}
	res, err := db.NamedExec(createNewTodoSchema, todo)
	if err != nil {
		HandleError(err, rw, t)
		return
	}
	todoId, err := res.LastInsertId()
	if err != nil {
		HandleError(err, rw, t)
		return
	}
	res, err = db.Exec(createNewUserTodoSchema, userId, todoId)
	if err != nil {
		HandleError(err, rw, t)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("TODO added successfully"))
}

func HandleError(err error, rw http.ResponseWriter, t Todos) {
	t.l.Panicln(err)
	rw.WriteHeader(http.StatusInternalServerError)
}
