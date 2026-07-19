package container

import (
	"fmt"

	"github.com/omnlgy/jadwalin/internal/config"
	"github.com/omnlgy/jadwalin/internal/controller"
	"github.com/omnlgy/jadwalin/internal/db"
	"github.com/omnlgy/jadwalin/internal/domain"
	"github.com/omnlgy/jadwalin/internal/models"
	"github.com/omnlgy/jadwalin/internal/provider"
	"github.com/omnlgy/jadwalin/internal/repository"
	"github.com/omnlgy/jadwalin/internal/service"
)

// Container holds all application dependencies
type Container struct {
	AuthController       *controller.Auth
	UserController       *controller.User
	TreatmentController  *controller.Treatment
	StaffSkillController *controller.StaffSkill
	BookingController    *controller.Booking
	AuthService          domain.AuthService
}

// InitializeContainer creates and initializes all application components
func InitializeContainer(cfg *config.Config) (*Container, error) {
	posgreDb, err := db.NewPostgresDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate tables on startup
	if err := posgreDb.AutoMigrate(&models.User{}, &models.Treatment{}, &models.Booking{}, &models.StaffSkill{}); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	rDb := db.NewRedisClient(cfg)

	// Initialize repositories
	userRepo := repository.NewUserRepository(posgreDb)
	authRepo := repository.NewAuthRepository(rDb)
	treatmentRepo := repository.NewTreatmentRepository(posgreDb)
	staffSkillRepo := repository.NewStaffSkillRepository(posgreDb)
	bookingRepo := repository.NewBookingRepository(posgreDb)

	// Initialize notification provider
	waProvider := provider.NewWhatsAppProvider(cfg.GOWA_URL, cfg.GOWA_DEVICE_ID)
	emailProvider := provider.NewEmailProvider(cfg.SMTP_HOST, cfg.SMTP_PORT, cfg.SMTP_USER, cfg.SMTP_PASS, cfg.SMTP_SENDER)

	// Initialize services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(authRepo)
	notificationService := service.NewNotificationService(waProvider, emailProvider)
	treatmentService := service.NewTreatmentService(treatmentRepo)
	staffSkillService := service.NewStaffSkillService(staffSkillRepo)
	bookingService := service.NewBookingService(bookingRepo, userRepo, treatmentRepo, staffSkillRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService, userService, notificationService)
	userController := controller.NewUserController(userService, authService, notificationService)
	treatmentController := controller.NewTreatmentController(treatmentService)
	staffSkillController := controller.NewStaffSkillController(staffSkillService)
	bookingController := controller.NewBookingController(bookingService, treatmentService, notificationService)

	return &Container{
		AuthController:       authController,
		UserController:       userController,
		TreatmentController:  treatmentController,
		StaffSkillController: staffSkillController,
		BookingController:    bookingController,
		AuthService:          authService,
	}, nil
}

// Close gracefully closes any open connections in the container
func (c *Container) Close() error {
	// Add logic to close database connections, Redis, etc.
	// For GORM, you might need to get the underlying *sql.DB and close it.
	// Example:
	// if sqlDB, err := c.AuthRepository.DB.DB(); err == nil {
	// 	sqlDB.Close()
	// }
	return nil
}

// TODO: Need to add a way to get the *gorm.DB instance from one of the repos to close it.
// For now, defer sqlDb.Close() in main.go is sufficient.
