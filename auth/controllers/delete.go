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

	// userService := service.NewUserService()
	err := a.UserService.Delete(username)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Deleted user successfully."))
}
