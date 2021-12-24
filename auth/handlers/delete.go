package handlers

import (
	"fmt"
	auth_db "go-todo/auth/db"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var deleteUserSchema = `DELETE FROM users where username=?`

func (a Auth) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	res, err := auth_db.GetDb().Exec(deleteUserSchema, username)
	if err != nil {
		a.l.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	c, err := res.RowsAffected()
	if err != nil {
		a.l.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if c != 1 {
		m := fmt.Sprintf("Could not find user with username - %s", username)
		a.l.Println(m)
		rw.Write([]byte(m))
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Deleted user successfully."))
}
