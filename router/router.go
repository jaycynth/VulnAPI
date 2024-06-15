package router

import (
	"net/http"

	"github.com/security-testing-api/controller"
	"github.com/security-testing-api/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(userController *controller.UserController, kycController *controller.KYCController) *gin.Engine {
	service := gin.Default()

	service.Use(middleware.CORSMiddleware())

	service.GET("", func(context *gin.Context) {
		context.JSON(http.StatusOK, "welcome home")
	})

	service.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r := service.Group("/api")
	r.Use(middleware.CORSMiddleware())

	r.POST("/login", userController.Login)
	r.POST("/register", userController.Create)

	userProtectedRoutes := r.Group("/users")
	userProtectedRoutes.Use(middleware.AuthMiddleware())
	{
		// userProtectedRoutes.GET("/:userId", userController.FindById)
		userProtectedRoutes.GET("", userController.FindAll)
	}

	kycProtectedRoutes := r.Group("/kyc")
	kycProtectedRoutes.Use(middleware.AuthMiddleware())
	{
		kycProtectedRoutes.POST("", kycController.Save)
	}

	return service
}
