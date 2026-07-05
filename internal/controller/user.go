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
	u.userService.RegisterEmployee(newUser)

	// todo: return user data
	// send otp to verify phone number
}
