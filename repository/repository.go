package repository

import "github.com/Falcer/go-login/model"

type (
	// UserRepository interface
	UserRepository interface {
		GetAllUser() ([]model.User, error)
		GetUserByUsername(username string) (*model.User, error)
		AddUser(user model.User) error
	}
)
