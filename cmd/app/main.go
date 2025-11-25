package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pndwrzk/go-article/config"
	"github.com/pndwrzk/go-article/internal/middleware"
	"github.com/pndwrzk/go-article/pkg/database"
	"github.com/pndwrzk/go-article/routes"
)

func main() {

	config.LoadConfig()
	database.ConnectPostgres()
	database.Migrate()

	router := gin.Default()
	router.Use(middleware.InjectHost())

	routes.RegisterRoutes(router)
	router.Static("/uploads", "./uploads")

	port := config.AppConfig.Port
	router.Run(":" + port)
}
