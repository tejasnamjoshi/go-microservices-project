package handlers

import (
	"encoding/json"
	"go-todo/auth/data"
	auth_db "go-todo/auth/db"
	"go-todo/auth/repository"
	"go-todo/auth/utils"
	"net/http"
)

type LoginResp struct {
	Token string `json:"token"`
}

func (a Auth) AddUser(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user = data.User{}
	user.FromJSON(r.Body)
	err := user.GeneratePassword()
	if err != nil {
		utils.CreateHttpError(rw, http.StatusInternalServerError, "Error creating user.")
		return
	}

	err = a.v.Struct(user)
	if err != nil {
		a.l.Error(err)
		utils.CreateHttpError(rw, http.StatusBadRequest, "Error creating user.")
		return
	}
	db := auth_db.GetDb()

	_, err = repository.StoreUser(db, user, a.l)
	if err != nil {
		utils.CreateHttpError(rw, http.StatusInternalServerError, "Error creating user.")
		return
	}
	if err != nil {
		utils.CreateHttpError(rw, http.StatusInternalServerError, "Error creating user.")
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("New User created successfully"))
}

func (a Auth) Login(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	user := data.User{}

	user.FromJSON(r.Body)
	err := a.v.Struct(user)
	if err != nil {
		a.l.Error(err)
		utils.CreateHttpError(rw, http.StatusBadRequest, "Error logging in.")
		return
	}

	db := auth_db.GetDb()

	token, err := repository.AuthenticateUser(db, user, a.l)
	if err != nil {
		utils.CreateHttpError(rw, http.StatusUnauthorized, "Error logging in.")
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	e := json.NewEncoder(rw)
	e.Encode(LoginResp{Token: token})
}
