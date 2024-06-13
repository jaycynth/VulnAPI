package repository

import (
	"errors"

	"github.com/security-testing-api/data/request"
	"github.com/security-testing-api/helper"
	"github.com/security-testing-api/model"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	Db *gorm.DB
}

func NewUserRepositoryImpl(Db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{Db: Db}
}

func (t UserRepositoryImpl) Save(user model.User) (model.User, error) {
	switch {
	case user.Username == "":
		return model.User{}, helper.ErrorFormatter(helper.ErrMissingData, "Username")
	case user.Password == "":
		return model.User{}, helper.ErrorFormatter(helper.ErrMissingData, "Password")
	case user.Email == "":
		return model.User{}, helper.ErrorFormatter(helper.ErrMissingData, "Email")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, helper.ErrorFormatter(err, "Error encrypting password")
	}
	user.Password = string(hashedPassword)

	err = t.Db.Create(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return model.User{}, helper.ErrDuplicateData
		}
		return model.User{}, helper.ErrorFormatter(err, "Error saving user")
	}

	return user, nil

}

func (t UserRepositoryImpl) Login(username string, password string) (model.User, error) {
	switch {
	case username == "":
		return model.User{}, helper.ErrorFormatter(helper.ErrMissingData, "Username")
	case password == "":
		return model.User{}, helper.ErrorFormatter(helper.ErrMissingData, "Password")
	}

	var user model.User

	err := t.Db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, helper.ErrUserNotFound
		}
		return model.User{}, helper.ErrorFormatter(err, "Error fetching user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return model.User{}, helper.ErrInvalidData
	}

	return user, nil
}

func (t UserRepositoryImpl) FindAll() ([]model.User, error) {

	var users []model.User

	err := t.Db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil

}

func (t UserRepositoryImpl) Update(User model.User) {
	var updateTag = request.UpdateUserRequest{
		Id:       User.Id,
		Username: User.Username,
	}
	result := t.Db.Model(&User).Updates(updateTag)
	helper.ErrorPanic(result.Error)
}

func (t UserRepositoryImpl) Delete(UserId int) {
	var User model.User
	result := t.Db.Where("id = ?", UserId).Delete(&User)
	helper.ErrorPanic(result.Error)
}

func (t UserRepositoryImpl) FindById(UserId int) (model.User, error) {
	var user model.User
	result := t.Db.Find(&user, UserId)
	if result != nil {
		return user, nil
	} else {
		return user, errors.New("user is not found")
	}
}
