package router

import (
	"net/http"

	"github.com/security-testing-api/controller"
	"github.com/security-testing-api/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(userController *controller.UserController) *gin.Engine {
	service := gin.Default()

	service.GET("", func(context *gin.Context) {
		context.JSON(http.StatusOK, "welcome home")
	})

	service.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r := service.Group("/api")

	r.POST("/login", userController.Login)
	r.POST("/register", userController.Create)

	protectedRoutes := r.Group("/").Use(middleware.TokenAuthMiddleware())
	{
		protectedRoutes.PUT("/users/:userId", userController.Update)
		protectedRoutes.DELETE("/users/:userId", userController.Delete)
		protectedRoutes.GET("/users/:userId", userController.FindById)
		protectedRoutes.GET("/users", userController.FindAll)
	}

	return service
}
