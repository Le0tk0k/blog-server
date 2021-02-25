package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Le0tk0k/blog-server/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository はUserRepositoryを返す
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

// FindByEmail はemailを持つuserを取得する
func (u *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := u.db.Get(&user, "SELECT email, password FROM users WHERE email = ?", email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return nil, fmt.Errorf("FindByEmail: cannot find user: %w", model.ErrUserNotFound)
		}
		return nil, fmt.Errorf("FindByEmail: cannot find user: %w", err)
	}
	return &user, nil
}
