package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	ConnectDB()

	r := gin.Default()

	r.POST("/register", Register)
	r.POST("/login", Login)

	protected := r.Group("/api")
	protected.Use(AuthMiddleware())
	protected.GET("/me", Protected)

	r.Run(":8082")
}
