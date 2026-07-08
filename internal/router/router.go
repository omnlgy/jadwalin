package router

import (
	"github.com/gin-gonic/gin"
	"github.com/omnlgy/jadwalin/internal/controller"
	"github.com/omnlgy/jadwalin/internal/middleware"
)

func AuthRoutes(router *gin.Engine, controller controller.Auth) {
	auth := router.Group("/api/auth")

	auth.POST("/register-otp", controller.RegisterOTP)
	auth.POST("/login", controller.Login)
	auth.POST("/login-verify", controller.LoginVerify)
}

func TreatmentRoutes(router *gin.Engine, controller controller.Treatment) {
	treatment := router.Group("/api/treatment")

	treatment.POST("/", middleware.AuthMiddleware(), middleware.RequireRole("admin"), controller.CreateTreatment)
	treatment.GET("/list", controller.ListTreatments)
	treatment.GET("/:id", controller.GetTreatment)
	treatment.PUT("/:id", middleware.AuthMiddleware(), controller.UpdateTreatment)
	treatment.DELETE("/:id", middleware.AuthMiddleware(), middleware.RequireRole("admin"), controller.DeleteTreatment)
}

func UserRoutes(router *gin.Engine, controller controller.User) {
	user := router.Group("/api/user")

	user.POST("/register-staff", middleware.AuthMiddleware(), middleware.RequireRole("admin"), controller.RegisterStaff)
	user.POST("/verify", controller.VerifyUser)
	user.GET("/list", controller.ListUsers)
	user.PUT("/:id", middleware.AuthMiddleware(), controller.UpdateUser)
	user.DELETE("/:id", middleware.AuthMiddleware(), middleware.RequireRole("admin"), controller.DeleteUser)
}
