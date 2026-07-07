package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omnlgy/jadwalin/internal/domain"
	"github.com/omnlgy/jadwalin/internal/dto"
)

type Auth struct {
	authService domain.AuthService
	userService domain.UserService
}

func NewAuthController(authService domain.AuthService, userService domain.UserService) *Auth {
	return &Auth{
		authService: authService,
		userService: userService,
	}
}

// RegisterOTP godoc
// @Summary Register OTP for a user
// @Description Sends an OTP to the user's registered phone number for verification.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterOTPRequest true "Register OTP Request"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.BadRequestResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/auth/register-otp [post]
func (c *Auth) RegisterOTP(ctx *gin.Context) {
	var body dto.RegisterOTPRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		abortWithBadRequest(ctx, err)
		return
	}
	userId, err := uuid.Parse(body.UserID)
	if err != nil {
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: "bad request",
			Errors: []dto.FieldError{
				{
					Field:   "user_id",
					Message: "invalid user id",
				},
			},
		})
		return
	}

	user, err := c.userService.GetByID(userId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
				Code:    400,
				Message: "user not found",
			})
		} else {
			ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return
	}

	if user.PhoneNumber != body.Phone {
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: "bad request",
			Errors: []dto.FieldError{
				{
					Field:   "phone",
					Message: "phone number is invalid",
				},
			},
		})
		return
	}

	key := "register-otp:" + user.PhoneNumber
	err = c.authService.GenerateOTP(ctx.Request.Context(), key)
	if err != nil {
		ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, dto.SuccessResponse{
		Code:    200,
		Message: "otp sent",
	})
}
