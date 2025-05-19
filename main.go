package main

import (
	"api-sec-go/config"
	"api-sec-go/routes"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "api-sec-go/docs"
)

func main() {
	config.ConnectDB()

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.SetupRoutes(r)
	r.Run(":8080")
}
