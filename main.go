package main

import (
	"fiscus/database"
	"fiscus/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	database.MigrateTables()
	db := database.GetDB()
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // replace with specific origins in production when using AllowCredentials
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	RegisterAuthRoutes(server)

	protected := server.Group("/")
	protected.Use(middleware.Authorize)
	RegisterAccountRoutes(db, protected)
	RegisterCategoryRoutes(db, protected)
	RegisterTrxRoutes(db, protected)

	server.Run(":8080")
}
