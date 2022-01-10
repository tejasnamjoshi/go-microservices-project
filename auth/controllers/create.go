package controllers

import (
	"go-todo/auth/entities"
	"net/http"
)

type LoginResp struct {
	Token string `json:"token"`
}

func (a Auth) AddUser(rw http.ResponseWriter, r *http.Request) {
	// Format Input
	defer r.Body.Close()
	var user = entities.User{}
	user.FromJSON(r.Body)

	// Validate Input
	err := a.UserService.Validate(&user)
	if err != nil {
		a.Logger.Error(err.Error())
		a.Response.SendErrorResponse(rw, http.StatusBadRequest, err.Error())
		return
	}

	// Invoke Logic
	err = a.UserService.Create(&user)
	if err != nil {
		a.Response.SendErrorResponse(rw, http.StatusInternalServerError, "Error creating user.")
		return
	}

	// Send Response
	a.Response.SendSuccessResponse(rw, "New User created successfully")
}

func (a Auth) Login(rw http.ResponseWriter, r *http.Request) {
	// Format Input
	defer r.Body.Close()
	user := entities.User{}
	user.FromJSON(r.Body)

	// Validate Input
	err := a.UserService.Validate(&user)
	if err != nil {
		a.Logger.Error(err.Error())
		a.Response.SendErrorResponse(rw, http.StatusBadRequest, err.Error())
		return
	}

	// Invoke Logic
	token, err := a.UserService.Login(&user)
	if err != nil {
		a.Logger.Error(err.Error())
		a.Response.SendErrorResponse(rw, http.StatusUnauthorized, "Error logging in.")
		return
	}

	// Send Response
	a.Response.SendSuccessResponse(rw, LoginResp{Token: token})
}
