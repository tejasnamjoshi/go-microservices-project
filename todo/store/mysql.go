package store

import (
	"fmt"
	"go-todo/todo/logging"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db_client *sqlx.DB

func initDB(logger logging.Logger) {
	uri := fmt.Sprintf("%s:%s@tcp(%s)/go-todo", os.Getenv("MYSQL_ROOT_USERNAME"), os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_HOST"))
	db, err := sqlx.Connect("mysql", uri)
	if err != nil {
		logger.Fatal(err.Error())
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		logger.Fatal(err.Error())
		panic(err)
	}

	db_client = db
}

func GetDb(logger logging.Logger) *sqlx.DB {
	if db_client == nil {
		initDB(logger)
	}

	return db_client
}
