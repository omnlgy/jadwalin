package router

import (
	"github.com/gin-gonic/gin"
	"github.com/omnlgy/jadwalin/internal/controller"
)

func AuthRoutes(router *gin.Engine, controller controller.Auth) {
	auth := router.Group("/api/auth")

	auth.POST("/register-otp", controller.RegisterOTP)
	auth.POST("/login", controller.Login)
	auth.POST("/login-verify", controller.LoginVerify)
}

func UserRoutes(router *gin.Engine, controller controller.User) {
	user := router.Group("/api/user")

	user.POST("/register-employee", controller.RegisterEmployee)
	user.POST("/verify", controller.VerifyUser)
	user.GET("/list", controller.ListUsers)
}
