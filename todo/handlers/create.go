package handlers

import (
	"go-todo/todo/entities"
	"net/http"
)

func (t Todos) CreateNewTodo(rw http.ResponseWriter, r *http.Request) {
	var todo = entities.Todo{}
	ctx := r.Context()
	userId := ctx.Value("userId").(int)

	defer r.Body.Close()
	err := todo.FromJSON(r.Body)
	if err != nil {
		HandleError(err, rw, t)
		return
	}

	err = t.todoRepository.Create(&todo, userId)
	if err != nil {
		HandleError(err, rw, t)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("TODO added successfully"))
}

func HandleError(err error, rw http.ResponseWriter, t Todos) {
	t.l.Error(err.Error())
	rw.WriteHeader(http.StatusInternalServerError)
}
