package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/omnlgy/jadwalin/internal/db"
)

func HealthRoutes(router *gin.Engine) {
	router.GET("/api/health", func(ctx *gin.Context) {
		dbStatus := "connected"
		redisStatus := "connected"
		httpStatus := http.StatusOK

		// Check PostgreSQL
		sqlDB, err := db.DB.DB()
		if err != nil {
			dbStatus = "unavailable"
			httpStatus = http.StatusServiceUnavailable
		} else if err := sqlDB.Ping(); err != nil {
			dbStatus = "disconnected"
			httpStatus = http.StatusServiceUnavailable
		}

		// Check Redis
		if _, err := db.Redis.Ping(context.Background()).Result(); err != nil {
			redisStatus = "disconnected"
			httpStatus = http.StatusServiceUnavailable
		}

		ctx.JSON(httpStatus, gin.H{
			"status":   "ok",
			"database": dbStatus,
			"redis":    redisStatus,
		})
	})
}
