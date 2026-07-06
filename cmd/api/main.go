package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/omnlgy/jadwalin/internal/config"
	"github.com/omnlgy/jadwalin/internal/controller"
	"github.com/omnlgy/jadwalin/internal/db"
	"github.com/omnlgy/jadwalin/internal/repository"
	"github.com/omnlgy/jadwalin/internal/router"
	"github.com/omnlgy/jadwalin/internal/service"
)

func main() {
	fmt.Println("Initializing server...")

	cfg := config.Load()

	posgreDb, err := db.NewPostgresDB(cfg)
	if err != nil {
		fmt.Println("Failed to connect to database")
		return
	}

	if sqlDb, err := posgreDb.DB(); err != nil {
		fmt.Println("Failed to get database connection")
		return
	} else {
		defer sqlDb.Close()
	}

	rDb := db.NewRedisClient(cfg)
	defer rDb.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(posgreDb)
	authRepo := repository.NewAuthRepository(rDb)

	// Initialize services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(authRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService, userService)
	userController := controller.NewUserController(userService, authService)

	server := gin.New()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())

	router.AuthRoutes(server, *authController)
	router.UserRoutes(server, *userController)

	fmt.Println("Starting server on port " + cfg.APP_PORT)
	if err := server.Run(":" + cfg.APP_PORT); err != nil {
		fmt.Println("Failed to start server")
		return
	}
}
