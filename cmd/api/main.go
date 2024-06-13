package main

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/security-testing-api/config"
	"github.com/security-testing-api/controller"
	"github.com/security-testing-api/helper"
	"github.com/security-testing-api/model"
	"github.com/security-testing-api/repository"
	"github.com/security-testing-api/router"
	"github.com/security-testing-api/service"
)

func main() {

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		helper.ErrorPanic(err)
	}

	db, err := config.DatabaseConnection(cfg.Database)
	if err != nil {
		helper.ErrorPanic(err)
	}

	validate := validator.New()

	db.Table("users").AutoMigrate(&model.User{})

	userRepository := repository.NewUserRepositoryImpl(db)

	userService := service.NewUserServiceImpl(userRepository, validate)

	userController := controller.NewUserController(userService)

	routes := router.NewRouter(userController)

	server := &http.Server{
		Addr:           ":8888",
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err = server.ListenAndServe()
	helper.ErrorPanic(err)

}
