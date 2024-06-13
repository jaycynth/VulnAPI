package service

import (
	"github.com/security-testing-api/data/request"
	"github.com/security-testing-api/data/response"
)

type UserService interface {
	Create(User request.CreateUserRequest)
	Update(User request.UpdateUserRequest)
	Delete(UserId int)
	FindById(UserId int) response.UserResponse
	FindAll() ([]response.UserResponse, error)
	Login(Username string, Password string) (response.UserResponse, error)
}
