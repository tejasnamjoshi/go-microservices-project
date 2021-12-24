package handlers

import (
	"encoding/json"
	"go-todo/auth/data"
	auth_db "go-todo/auth/db"
	"net/http"
)

type AuthRequest struct {
	Token string `json:"token"`
}

var selectAllSchema = `SELECT * FROM users`
var getByUserNameSchema = "SELECT * FROM users where username=?"

func (a Auth) GetUsers(rw http.ResponseWriter, r *http.Request) {
	users := data.Users{}
	err := auth_db.GetDb().Select(&users, selectAllSchema)
	if err != nil {
		a.l.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	users.ToJSON(rw)

	a.l.Printf("Users Fetched")
}

func (a Auth) GetUserByUsername(username string) (data.User, error) {
	var user = data.User{}
	err := auth_db.GetDb().Get(&user, getByUserNameSchema, username)
	if err != nil {
		a.l.Println(err)
		return user, err
	}

	return user, err
}

func (a Auth) GetUserAuthStatus(rw http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		a.l.Println("Token not available")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	ok, err := data.GetAuthorizationStatus("Bearer " + token)
	if err != nil {
		a.l.Println(err)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Unautorized access"))
		return
	}
	if ok {
		claims, err := data.ParseJWT(token)
		if err != nil {
			a.l.Println(err)
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unautorized access"))
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		e := json.NewEncoder(rw)
		e.Encode(claims)
		return
	}
	rw.Write([]byte("Unautorized access"))
	rw.WriteHeader(http.StatusUnauthorized)
}
