package handlers

import (
	"go-todo/auth/data"
	auth_db "go-todo/auth/db"
	"net/http"
)

var postSchema = `INSERT INTO users (username, password) VALUES (:username, :password)`

func (a Auth) AddUser(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user = data.User{}
	user.FromJSON(r.Body)

	res, err := auth_db.GetDb().NamedExec(postSchema, user)
	if err != nil {
		a.l.Panicln(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
	rows, err := res.RowsAffected()
	if rows == 0 {
		a.l.Panicln(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Add("Content-Type", "application/json")
	rw.Write([]byte("New User created successfully"))
}
