package handlers

import (
	"fmt"
	"go-todo/auth/data"
	auth_db "go-todo/auth/db"
	"net/http"
)

var addUserSchema = `INSERT INTO users (username, password) VALUES (:username, :password)`
var authUserSchema = `SELECT * FROM users WHERE username=?`

func (a Auth) AddUser(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user = data.User{}
	user.FromJSON(r.Body)
	err := user.GeneratePassword()
	if err != nil {
		a.l.Panicln(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error while logging in."))
		return
	}

	res, err := auth_db.GetDb().NamedExec(addUserSchema, user)
	if err != nil {
		a.l.Panicln(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
	rows, err := res.RowsAffected()
	if rows == 0 {
		a.l.Panicln(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("New User created successfully"))
}

func (a Auth) Login(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	validUsers := data.Users{}
	user := data.User{}

	user.FromJSON(r.Body)
	err := auth_db.GetDb().Select(&validUsers, authUserSchema, user.Username)
	if err != nil {
		a.l.Panicln(err)
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
		a.l.Panicln(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error logging in"))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Add("Content-Type", "application/json")
	resp := fmt.Sprintf("Valid user : %s", token)
	rw.Write([]byte(resp))
}
