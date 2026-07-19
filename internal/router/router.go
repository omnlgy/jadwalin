package router

import (
	"github.com/gin-gonic/gin"
	"github.com/omnlgy/jadwalin/internal/container"
	"github.com/omnlgy/jadwalin/internal/middleware"
)

func AuthRoutes(router *gin.Engine, cont *container.Container) {
	auth := router.Group("/api/auth")

	auth.POST("/register-otp", cont.AuthController.RegisterOTP)
	auth.POST("/login", cont.AuthController.Login)
	auth.POST("/login-verify", cont.AuthController.LoginVerify)
	auth.POST("/logout", middleware.AuthMiddleware(), cont.AuthController.Logout)
}

func StaffSkillRoutes(router *gin.Engine, cont *container.Container) {
	group := router.Group("/api/staff-skills")
	group.POST("/", middleware.AuthMiddleware(), middleware.RequireRole("admin"), cont.StaffSkillController.AssignSkill)
	group.GET("/list", cont.StaffSkillController.ListAll)
	group.GET("/:id", cont.StaffSkillController.GetStaffSkill)
	group.GET("/staff/:userId", cont.StaffSkillController.ListByStaff)
	group.GET("/treatment/:treatmentId", cont.StaffSkillController.ListByTreatment)
	group.DELETE("/:id", middleware.AuthMiddleware(), middleware.RequireRole("admin"), cont.StaffSkillController.UnassignSkill)
}

func TreatmentRoutes(router *gin.Engine, cont *container.Container) {
	treatment := router.Group("/api/treatment")

	treatment.POST("/", middleware.AuthMiddleware(), middleware.RequireRole("admin"), cont.TreatmentController.CreateTreatment)
	treatment.GET("/list", cont.TreatmentController.ListTreatments)
	treatment.GET("/:id", cont.TreatmentController.GetTreatment)
	treatment.PUT("/:id", middleware.AuthMiddleware(), cont.TreatmentController.UpdateTreatment)
	treatment.DELETE("/:id", middleware.AuthMiddleware(), middleware.RequireRole("admin"), cont.TreatmentController.DeleteTreatment)
}

func UserRoutes(router *gin.Engine, cont *container.Container) {
	user := router.Group("/api/user")

	user.POST("/register-staff", middleware.AuthMiddleware(), middleware.RequireRole("admin"), cont.UserController.RegisterStaff)
	user.POST("/verify", cont.UserController.VerifyUser)
	user.GET("/list", cont.UserController.ListUsers)
	user.PUT("/:id", middleware.AuthMiddleware(), cont.UserController.UpdateUser)
	user.DELETE("/:id", middleware.AuthMiddleware(), middleware.RequireRole("admin"), cont.UserController.DeleteUser)
	user.POST("/:id/photo", middleware.AuthMiddleware(), cont.UserController.UploadPhoto)
}

func BookingRoutes(router *gin.Engine, cont *container.Container) {
	group := router.Group("/api/booking")
	group.POST("/available-slots", cont.BookingController.GetAvailableSlots)
	group.POST("/", middleware.AuthMiddleware(), cont.BookingController.CreateBooking)
	group.GET("/:id", middleware.AuthMiddleware())
	group.GET("/user/:userId", middleware.AuthMiddleware(), cont.BookingController.GetByUserID)
}
