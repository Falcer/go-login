package repository

import "github.com/Falcer/go-login/model"

type (
	// UserRepository interface
	UserRepository interface {
		GetUserByUsername(username string) (*model.User, error)
		AddUser(user model.User) error
	}
)
