package repository

import (
	"fmt"

	"github.com/Le0tk0k/blog-server/config"
	"github.com/jmoiron/sqlx"
)

// NewDB はMySQLへ接続し、DBオブジェクトを返す
func NewDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL : %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping : %w", err)
	}

	return db, nil
}
