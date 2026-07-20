package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/omnlgy/jadwalin/internal/db"
	"github.com/omnlgy/jadwalin/internal/storage"
)

func HealthRoutes(router *gin.Engine) {
	router.GET("/api/health", func(ctx *gin.Context) {
		dbStatus := "connected"
		redisStatus := "connected"
		storageStatus := "connected"
		status := true

		// Check PostgreSQL
		sqlDB, err := db.DB.DB()
		if err != nil {
			dbStatus = "unavailable"
			status = false
		} else if err := sqlDB.Ping(); err != nil {
			dbStatus = "disconnected"
			status = false
		}

		// Check Redis
		if _, err := db.Redis.Ping(context.Background()).Result(); err != nil {
			redisStatus = "disconnected"
			status = false
		}

		if _, err := storage.MinioClient.ListBuckets(ctx); err != nil {
			storageStatus = "disconnected"
			status = false
		}

		httpStatus := http.StatusOK
		if !status {
			httpStatus = http.StatusServiceUnavailable
		}

		ctx.JSON(httpStatus, gin.H{
			"status":   "ok",
			"database": dbStatus,
			"redis":    redisStatus,
			"storage":  storageStatus,
		})
	})
}
