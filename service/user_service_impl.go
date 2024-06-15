package service

import (
	"github.com/security-testing-api/data/request"
	"github.com/security-testing-api/data/response"
	"github.com/security-testing-api/helper"
	"github.com/security-testing-api/model"
	"github.com/security-testing-api/repository"

	"github.com/go-playground/validator/v10"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	Validate       *validator.Validate
}

func NewUserServiceImpl(userRepository repository.UserRepository, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		Validate:       validate,
	}
}

func (t UserServiceImpl) Create(user request.CreateUserRequest) {
	err := t.Validate.Struct(user)
	helper.ErrorPanic(err)
	userModel := model.User{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}
	t.UserRepository.Save(userModel)
}

func (t UserServiceImpl) FindAll() ([]response.UserResponse, error) {
	result, err := t.UserRepository.FindAll()
	if err != nil {
		return nil, err
	}

	var users []response.UserResponse
	for _, value := range result {
		user := response.UserResponse{
			Id:       value.Id,
			Username: value.Username,
			Email:    value.Email,
		}
		users = append(users, user)
	}
	return users, nil
}

func (t UserServiceImpl) Login(Username string, Passwrod string) (response.UserResponse, error) {
	result, err := t.UserRepository.Login(Username, Passwrod)
	if err != nil {
		return response.UserResponse{}, err
	}

	user := response.UserResponse{
		Id:       result.Id,
		Username: result.Username,
		Email:    result.Email,
	}

	return user, nil
}
