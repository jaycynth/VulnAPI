package repository

import (
	"github.com/security-testing-api/model"
)

type UserRepository interface {
	Save(user *model.User) (userRes *model.User, err error)
	FindAll() (users []*model.User, err error)
	Login(Username string, Password string) (user *model.User, err error)
}
