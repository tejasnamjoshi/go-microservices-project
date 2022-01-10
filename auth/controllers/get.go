package controllers

import (
	"net/http"
)

type AuthRequest struct {
	Token string `json:"token"`
}

func (a Auth) GetUsers(rw http.ResponseWriter, r *http.Request) {
	// invoke logic
	users, err := a.UserService.GetAll()
	if err != nil {
		a.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send Response
	a.Response.SendSuccessResponse(rw, users)
}

func (a Auth) DecodeToken(rw http.ResponseWriter, r *http.Request) {
	// Validate input
	token := r.URL.Query().Get("token")
	if token == "" {
		a.Logger.Error("Token not available")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Invoke logic
	claims, err := a.JwtService.GetAuthorizationData(token)
	if err != nil {
		a.Logger.Error(err.Error())
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Unauthorized access"))
		return
	}

	// Send Response
	a.Response.SendSuccessResponse(rw, claims)
}
