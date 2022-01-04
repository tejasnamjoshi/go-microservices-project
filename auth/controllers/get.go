package controllers

import (
	"encoding/json"
	"net/http"
)

type AuthRequest struct {
	Token string `json:"token"`
}

func (a Auth) GetUsers(rw http.ResponseWriter, r *http.Request) {
	users, err := userService.GetAll()
	if err != nil {
		a.l.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	users.ToJSON(rw)

	a.l.Info("Users Fetched")
}

func (a Auth) DecodeToken(rw http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		a.l.Error("Token not available")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	claims, err := jwtService.GetAuthorizationData(token)
	if err != nil {
		a.l.Error(err)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Unautorized access"))
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	e := json.NewEncoder(rw)
	e.Encode(claims)
}
