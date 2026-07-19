package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/omnlgy/jadwalin/docs" // Import generated docs

	"github.com/omnlgy/jadwalin/internal/config"
	"github.com/omnlgy/jadwalin/internal/container"
	"github.com/omnlgy/jadwalin/internal/router"
)

func init() {
	time.Local = time.UTC
}

// @title Jadwalin API
// @version 1.0
// @description This is a sample server for a scheduling application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	fmt.Println("Initializing server...")

	cfg := config.Load()

	cont, err := container.InitializeContainer(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize container: %v\n", err)
		return
	}

	// TODO: Move posgreDb.DB().Close() and rDb.Close() into container.Close() after refactoring.
	// For now, assume posgreDb is managed internally by the container and its underlying *sql.DB is closed elsewhere.
	// Similarly for Redis.

	server := gin.New()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())

	// Serve static files
	server.Static("/uploads", "./public/uploads")

	router.AuthRoutes(server, cont)
	router.UserRoutes(server, cont)
	router.TreatmentRoutes(server, cont)
	router.StaffSkillRoutes(server, cont)
	router.BookingRoutes(server, cont)

	// Add Swagger UI
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("Starting server on port " + cfg.APP_PORT)
	if err := server.Run(":" + cfg.APP_PORT); err != nil {
		fmt.Println("Failed to start server")
		return
	}
}
