package repository

import (
	"github.com/security-testing-api/model"
)

type UserRepository interface {
	Save(User model.User) (UserRes model.User, err error)
	FindAll() (Users []model.User, err error)
	Login(Username string, Password string) (User model.User, err error)
}
