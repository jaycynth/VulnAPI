package service

import (
	"github.com/security-testing-api/data/request"
	"github.com/security-testing-api/data/response"
)

type UserService interface {
	Create(User request.CreateUserRequest)
	FindAll() ([]response.UserResponse, error)
	Login(Username string, Password string) (response.UserResponse, error)
}
