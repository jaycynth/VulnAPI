package repository

import (
	"errors"
	"strings"

	"github.com/security-testing-api/helper"
	"github.com/security-testing-api/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	Db     *gorm.DB
	Hasher helper.BcryptHasher
}

func NewUserRepositoryImpl(Db *gorm.DB, hasher helper.Hasher) UserRepository {
	return &UserRepositoryImpl{Db: Db, Hasher: helper.BcryptHasher{}}
}

func (t *UserRepositoryImpl) Save(user *model.User) (*model.User, error) {
	switch {
	case user.Username == "":
		return nil, errors.New("missing username")
	case user.Password == "":
		return nil, errors.New("missing password")
	case user.Email == "":
		return nil, errors.New("missing emai")
	}

	hashedPassword, err := t.Hasher.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error encrypting password")
	}
	user.Password = string(hashedPassword)

	err = t.Db.Create(&user).Error
	switch {
	case err == nil:
	case strings.Contains(strings.ToLower(err.Error()), "duplicate"):
		return nil, errors.New("duplicate")
	default:
		return nil, errors.New("error saving user")
	}

	return user, nil

}

func (t *UserRepositoryImpl) Login(username string, password string) (*model.User, error) {
	switch {
	case username == "":
		return nil, errors.New("missing username")
	case password == "":
		return nil, errors.New("missing password")
	}

	var user model.User

	err := t.Db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrUserNotFound
		}
		return nil, helper.ErrorFormatter(err, "Error fetching user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, helper.ErrInvalidData
	}

	return &user, nil
}

func (t *UserRepositoryImpl) FindAll() ([]*model.User, error) {

	var users []*model.User

	err := t.Db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil

}
