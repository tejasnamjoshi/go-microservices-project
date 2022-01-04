package handlers

import (
	auth_db "go-todo/auth/db"
	"go-todo/auth/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a Auth) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	db := auth_db.GetDb()
	err := repository.DeleteUser(db, username, a.l)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Deleted user successfully."))
}
