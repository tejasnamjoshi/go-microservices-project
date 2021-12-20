package todo_db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db_client *sqlx.DB

func initDB() {
	db, err := sqlx.Connect("mysql", "root:root@tcp(localhost:3306)/go-todo")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db_client = db
}

func GetDb() *sqlx.DB {
	if db_client == nil {
		initDB()
	}

	return db_client
}
