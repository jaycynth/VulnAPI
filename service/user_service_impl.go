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

func (t UserServiceImpl) Update(user request.UpdateUserRequest) {
	userData, err := t.UserRepository.FindById(user.Id)
	helper.ErrorPanic(err)
	userData.Username = user.Username
	t.UserRepository.Update(userData)
}

func (t UserServiceImpl) Delete(userId int) {
	t.UserRepository.Delete(userId)
}

func (t UserServiceImpl) FindById(userId int) response.UserResponse {
	userData, err := t.UserRepository.FindById(userId)
	helper.ErrorPanic(err)

	userResponse := response.UserResponse{
		Id:       userData.Id,
		Username: userData.Username,
	}
	return userResponse
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
	}

	return user, nil
}
