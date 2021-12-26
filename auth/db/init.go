package auth_db

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db_client *sqlx.DB

func initDB() {
	uri := fmt.Sprintf("%s:%s@tcp(%s)/go-todo", os.Getenv("MYSQL_ROOT_USERNAME"), os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_HOST"))
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
