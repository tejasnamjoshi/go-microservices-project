package auth_db

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db_client *sqlx.DB

func initDB() {
	uri := fmt.Sprintf("root:%s@tcp(localhost:3306)/go-todo", os.Getenv("ROOT_PASSWORD"))
	db, err := sqlx.Connect("mysql", uri)
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
