package service

import (
	"fmt"

	"github.com/Le0tk0k/blog-server/model"
	"github.com/Le0tk0k/blog-server/repository"
)

type UserService interface {
	GetUser(email string) (*model.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

// NewUserService はUserServiceを返す
func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

// GetUser はemailを持つ記事を取得する
func (u *userService) GetUser(email string) (*model.User, error) {
	user, err := u.userRepository.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("GetUser: cannot get user: %w", err)
	}
	return user, nil
}
