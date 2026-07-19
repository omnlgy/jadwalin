package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/omnlgy/jadwalin/internal/domain"
	"github.com/omnlgy/jadwalin/internal/dto"
	"github.com/omnlgy/jadwalin/utils"
)

func AuthMiddleware(authService domain.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.AbortWithStatusJSON(401, dto.UnauthorizedResponse{
				Code:    401,
				Message: "missing authorization header",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			ctx.AbortWithStatusJSON(401, dto.UnauthorizedResponse{
				Code:    401,
				Message: "invalid authorization format",
			})
			return
		}

		isBlacklisted, err := authService.IsBlacklisted(ctx, tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
				Code:    500,
				Message: "failed to check token blacklist",
			})
			return
		}
		if isBlacklisted {
			ctx.AbortWithStatusJSON(401, dto.UnauthorizedResponse{
				Code:    401,
				Message: "token is blacklisted",
			})
			return
		}

		claim, err := utils.ValidateJWT(tokenString)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidJWT) {
				ctx.AbortWithStatusJSON(401, dto.UnauthorizedResponse{
					Code:    401,
					Message: "invalid token",
				})
				return
			}
			ctx.AbortWithStatusJSON(401, dto.UnauthorizedResponse{
				Code:    401,
				Message: err.Error(),
			})
			return
		}

		ctx.Set("userId", claim.UserID)
		ctx.Set("phoneNumber", claim.PhoneNumber)
		ctx.Set("role", claim.Role)
		ctx.Set("tokenString", tokenString)
		ctx.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.GetString("role")
		for _, r := range roles {
			if r == role {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(403, dto.ForbiddenResponse{
			Code:    403,
			Message: "Access Forbidden",
		})
	}
}
