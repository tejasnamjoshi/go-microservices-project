package handlers

import (
	"go-todo/auth/data"
	auth_db "go-todo/auth/db"
	"net/http"
)

var selectAllSchema = `SELECT * FROM users`
var getByUserNameSchema = "SELECT * FROM users where username=?"

func (a Auth) GetUsers(rw http.ResponseWriter, r *http.Request) {
	users := data.Users{}
	err := auth_db.GetDb().Select(&users, selectAllSchema)
	if err != nil {
		a.l.Panicln(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Add("Content-Type", "application/json")
	users.ToJSON(rw)

	a.l.Printf("Users Fetched")
}

func (a Auth) GetUserByUsername(username string) (data.User, error) {
	var user = data.User{}
	err := auth_db.GetDb().Get(&user, getByUserNameSchema, username)
	if err != nil {
		a.l.Panicln(err)
		return user, err
	}

	return user, err
}
