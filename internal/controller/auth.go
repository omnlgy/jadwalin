package controller

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omnlgy/jadwalin/internal/domain"
	"github.com/omnlgy/jadwalin/internal/dto"
	"github.com/omnlgy/jadwalin/utils"
)

type Auth struct {
	authService         domain.AuthService
	userService         domain.UserService
	notificationService domain.NotificationService
}

func NewAuthController(authService domain.AuthService, userService domain.UserService, notificationService domain.NotificationService) *Auth {
	return &Auth{
		authService:         authService,
		userService:         userService,
		notificationService: notificationService,
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
	otpCode, err := c.authService.GenerateOTP(ctx.Request.Context(), key)
	if err != nil {
		ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	fmt.Println(otpCode)

	ctx.JSON(200, dto.SuccessResponse{
		Code:    200,
		Message: "otp sent",
	})
}

// Login godoc
// @Summary Login — request OTP
// @Description Sends an OTP to the registered phone number for login.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.BadRequestResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/auth/login [post]
func (c *Auth) Login(ctx *gin.Context) {
	var body dto.LoginRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		abortWithBadRequest(ctx, err)
		return
	}

	user, err := c.userService.GetByPhoneNumber(body.Phone)
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

	key := "login-otp:" + user.PhoneNumber
	otpCode, err := c.authService.GenerateOTP(ctx.Request.Context(), key)
	if err != nil {
		ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	if err := c.notificationService.SendOTPLoginWhatsApp(ctx.Request.Context(), user.PhoneNumber, otpCode); err != nil {
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

// LoginVerify godoc
// @Summary Login — verify OTP
// @Description Verifies the OTP for login and returns a JWT token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyUserRequest true "Verify OTP Request"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.BadRequestResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/auth/login-verify [post]
func (c *Auth) LoginVerify(ctx *gin.Context) {
	var body dto.VerifyUserRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		abortWithBadRequest(ctx, err)
		return
	}

	key := "login-otp:" + body.Phone
	err := c.authService.VerifyOTP(ctx.Request.Context(), key, body.OTP)
	if err != nil {
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: "invalid otp",
		})
		return
	}

	user, err := c.userService.GetByPhoneNumber(body.Phone)
	if err != nil {
		ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, dto.SuccessResponse{
		Code:    200,
		Message: "otp verified",
		Data: struct {
			Token string `json:"token"`
		}{
			Token: token,
		},
	})
}
