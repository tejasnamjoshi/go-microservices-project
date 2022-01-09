package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a Auth) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err := a.UserService.Delete(username)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	a.Response.SendSuccessResponse(rw, "Deleted user successfully")
}
