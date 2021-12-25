package handlers

import (
	"encoding/json"
	"go-todo/auth/data"
	auth_db "go-todo/auth/db"
	"net/http"
)

type LoginResp struct {
	Token string `json:"token"`
}

var addUserSchema = `INSERT INTO users (username, password) VALUES (:username, :password)`
var authUserSchema = `SELECT * FROM users WHERE username=?`

func (a Auth) AddUser(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user = data.User{}
	user.FromJSON(r.Body)

	err := a.v.Struct(user)
	if err != nil {
		a.l.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error creating user."))
		return
	}

	err = user.GeneratePassword()
	if err != nil {
		a.l.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error creating user."))
		return
	}

	res, err := auth_db.GetDb().NamedExec(addUserSchema, user)
	if err != nil {
		a.l.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rows, err := res.RowsAffected()
	if rows == 0 {
		a.l.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("New User created successfully"))
}

func (a Auth) Login(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	validUsers := data.Users{}
	user := data.User{}

	user.FromJSON(r.Body)
	err := a.v.Struct(user)
	if err != nil {
		a.l.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error logging in."))
		return
	}

	err = auth_db.GetDb().Select(&validUsers, authUserSchema, user.Username)
	if err != nil {
		a.l.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error logging in"))
		return
	}

	if len(validUsers) == 0 {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("User not found"))
		return
	}

	validUser := validUsers[0]
	isValid := user.ComparePassword(validUser.Password)

	if isValid == false {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Incorrect Credentials"))
		return
	}

	token, err := validUser.GetJWT()
	if err != nil {
		a.l.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error logging in"))
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	resp := LoginResp{Token: token}
	e := json.NewEncoder(rw)
	e.Encode(resp)
}
