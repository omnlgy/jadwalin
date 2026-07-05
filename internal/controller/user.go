package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/omnlgy/jadwalin/internal/domain"
	"github.com/omnlgy/jadwalin/internal/dto"
)

type User struct {
	userService domain.UserService
}

func NewUser(userService domain.UserService) *User {
	return &User{
		userService: userService,
	}
}

func (u *User) RegisterEmployee(ctx *gin.Context) {
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
	user, err := u.userService.RegisterEmployee(newUser)
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
