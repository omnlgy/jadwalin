package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/omnlgy/jadwalin/docs" // Import generated docs

	"github.com/omnlgy/jadwalin/internal/config"
	"github.com/omnlgy/jadwalin/internal/controller"
	"github.com/omnlgy/jadwalin/internal/db"
	"github.com/omnlgy/jadwalin/internal/models"
	"github.com/omnlgy/jadwalin/internal/provider"
	"github.com/omnlgy/jadwalin/internal/repository"
	"github.com/omnlgy/jadwalin/internal/router"
	"github.com/omnlgy/jadwalin/internal/service"
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

	// Auto-migrate tables on startup
	if err := posgreDb.AutoMigrate(&models.User{}, &models.Treatment{}, &models.Booking{}, &models.StaffSkill{}); err != nil {
		fmt.Println("Failed to auto-migrate database:", err)
		return
	}

	rDb := db.NewRedisClient(cfg)
	defer rDb.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(posgreDb)
	authRepo := repository.NewAuthRepository(rDb)
	treatmentRepo := repository.NewTreatmentRepository(posgreDb)
	staffSkillRepo := repository.NewStaffSkillRepository(posgreDb)
	bookingRepo := repository.NewBookingRepository(posgreDb)

	// Initialize notification provider
	waProvider := provider.NewWhatsAppProvider(cfg.GOWA_URL, cfg.GOWA_DEVICE_ID)

	// Initialize services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(authRepo)
	notificationService := service.NewNotificationService(waProvider, nil)
	treatmentService := service.NewTreatmentService(treatmentRepo)
	staffSkillService := service.NewStaffSkillService(staffSkillRepo)
	bookingService := service.NewBookingService(bookingRepo, userRepo, treatmentRepo, staffSkillRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService, userService, notificationService)
	userController := controller.NewUserController(userService, authService, notificationService)
	treatmentController := controller.NewTreatmentController(treatmentService)
	staffSkillController := controller.NewStaffSkillController(staffSkillService)
	bookingController := controller.NewBookingController(bookingService, treatmentService)

	server := gin.New()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())

	router.AuthRoutes(server, *authController)
	router.UserRoutes(server, *userController)
	router.TreatmentRoutes(server, *treatmentController)
	router.StaffSkillRoutes(server, *staffSkillController)
	router.BookingRoutes(server, *bookingController)

	// Add Swagger UI
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("Starting server on port " + cfg.APP_PORT)
	if err := server.Run(":" + cfg.APP_PORT); err != nil {
		fmt.Println("Failed to start server")
		return
	}
}
