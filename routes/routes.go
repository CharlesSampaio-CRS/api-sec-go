package routes

import (
	"api-sec-go/controllers"
	"api-sec-go/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	auth := r.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	auth.PUT("/user", controllers.Update)
}
