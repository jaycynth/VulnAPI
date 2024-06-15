package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/security-testing-api/data/request"
	"github.com/security-testing-api/data/response"
	"github.com/security-testing-api/helper"
	"github.com/security-testing-api/service"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{userService: service}
}

func (controller *UserController) Create(ctx *gin.Context) {

	createUserRequest := request.CreateUserRequest{}

	if ctx.Request.Method != http.MethodPost {
		ctx.JSON(http.StatusMethodNotAllowed, response.Response{
			Code:   http.StatusMethodNotAllowed,
			Status: "Method Not Allowed",
			Error:  "Only POST method is allowed",
		})
		return
	}

	if err := ctx.ShouldBindJSON(&createUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Error:  err.Error(),
		})
		return
	}

	controller.userService.Create(createUserRequest)

	webResponse := response.Response{
		Code:   200,
		Status: "Ok",
		Data:   nil,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *UserController) Login(ctx *gin.Context) {

	var loginRequest request.LoginRequest

	if ctx.Request.Method != http.MethodPost {
		ctx.JSON(http.StatusMethodNotAllowed, response.Response{
			Code:   http.StatusMethodNotAllowed,
			Status: "Method Not Allowed",
			Error:  "Only POST method is allowed",
		})
		return
	}

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Error:  err.Error(),
		})
		return
	}

	user, err := controller.userService.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		switch err {
		case helper.ErrInvalidData:
			ctx.JSON(http.StatusBadRequest, response.Response{
				Code:   http.StatusBadRequest,
				Status: "Bad Request",
			})
		case helper.ErrUserNotFound:
			ctx.JSON(http.StatusUnauthorized, response.Response{
				Code:   http.StatusUnauthorized,
				Status: "Unauthorized",
				Error:  err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, response.Response{
				Code:   http.StatusInternalServerError,
				Status: "Internal Server Error",
				Error:  err.Error(),
			})
		}
		return
	}

	token, err := helper.GenerateToken(user.Username, strconv.Itoa(user.Id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.ErrorFormatter(err, "Error generating token"))
		return
	}

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data: map[string]interface{}{
			"token": token,
			"user":  user,
		},
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *UserController) FindAll(ctx *gin.Context) {
	userResponse, err := controller.userService.FindAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:   http.StatusInternalServerError,
			Status: "Error",
			Error:  err.Error(),
		})
		return
	}

	webResponse := response.Response{
		Code:   200,
		Status: "Ok",
		Data:   userResponse,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)

}
