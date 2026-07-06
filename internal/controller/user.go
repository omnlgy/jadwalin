package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/omnlgy/jadwalin/internal/domain"
	"github.com/omnlgy/jadwalin/internal/dto"
)

type User struct {
	userService domain.UserService
	authService domain.AuthService
}

func NewUserController(userService domain.UserService, authService domain.AuthService) *User {
	return &User{
		userService: userService,
		authService: authService,
	}
}

func (c *User) RegisterEmployee(ctx *gin.Context) {
	var body dto.RegisterEmployeeRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		abortWithBadRequest(ctx, err)
		return
	}

	newUser := &domain.User{
		PhoneNumber: body.PhoneNumber,
		Email:       body.Email,
		FullName:    body.FullName,
		Photo:       body.Photo,
	}
	user, err := c.userService.RegisterEmployee(newUser)
	if err != nil {
		ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	dataUser := dto.User{
		Id:          user.ID.String(),
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		FullName:    user.FullName,
		Photo:       user.Photo,
		Role:        string(user.Role),
		Verified:    user.Verified,
	}

	ctx.JSON(201, dto.CreatedResponse{
		Code:    201,
		Message: "User created successfully",
		Data:    dataUser,
	})
}

func (c *User) VerifyUser(ctx *gin.Context) {
	var body dto.VerifyUserRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		abortWithBadRequest(ctx, err)
		return
	}

	key := "register-otp:" + body.Phone
	err := c.authService.VerifyOTP(ctx.Request.Context(), key, body.OTP)
	if err != nil {
		if err == domain.ErrInvalidOTP || err == domain.ErrNotFound {
			ctx.AbortWithStatusJSON(400, dto.InternalErrorResponse{
				Code:    400,
				Message: "Invalid OTP",
			})

		} else {
			ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return
	}

	ctx.JSON(200, dto.SuccessResponse{
		Code:    200,
		Message: "User verified successfully",
	})

}
