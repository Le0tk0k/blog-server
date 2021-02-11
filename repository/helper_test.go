package repository

import (
	"github.com/Le0tk0k/blog-server/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// NewTestDB はtest用DBオブジェクトを返す
func NewTestDB() *sqlx.DB {
	db, err := sqlx.Open("mysql", config.DSN())
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return db
}
