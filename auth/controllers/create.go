package controllers

import (
	"encoding/json"
	"go-todo/auth/entities"
	"go-todo/auth/utils"
	"net/http"
)

type LoginResp struct {
	Token string `json:"token"`
}

func (a Auth) AddUser(rw http.ResponseWriter, r *http.Request) {
	// Manipulate Input
	defer r.Body.Close()
	var user = entities.User{}
	user.FromJSON(r.Body)

	// Validate Input
	err := userService.Validate(&user)
	if err != nil {
		a.l.Error(err)
		utils.CreateHttpError(rw, http.StatusBadRequest, err.Error())
		return
	}

	// Invoke Logic
	err = userService.Create(&user)
	if err != nil {
		utils.CreateHttpError(rw, http.StatusInternalServerError, "Error creating user.")
		return
	}

	// Send Response
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("New User created successfully"))
}

func (a Auth) Login(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	user := entities.User{}

	user.FromJSON(r.Body)
	err := userService.Validate(&user)
	if err != nil {
		a.l.Error(err)
		utils.CreateHttpError(rw, http.StatusBadRequest, err.Error())
		return
	}

	token, err := userService.Login(&user)
	if err != nil {
		a.l.Error(err)
		utils.CreateHttpError(rw, http.StatusUnauthorized, "Error logging in.")
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	e := json.NewEncoder(rw)
	e.Encode(LoginResp{Token: token})
}
