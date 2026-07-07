package controller

import (
	"math"
	"strconv"

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

// RegisterEmployee godoc
// @Summary Register a new employee
// @Description Registers a new employee with the provided details.
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.RegisterEmployeeRequest true "Register Employee Request"
// @Success 201 {object} dto.CreatedResponse
// @Failure 400 {object} dto.BadRequestResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/user/register-employee [post]
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
		if err == domain.ErrConflict {
			ctx.AbortWithStatusJSON(409, dto.BadRequestResponse{
				Code:    409,
				Message: "User already exists",
			})
			return
		}
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

// VerifyUser godoc
// @Summary Verify user with OTP
// @Description Verifies a user by checking the provided OTP against the stored one.
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.VerifyUserRequest true "Verify User Request"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.InternalErrorResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/user/verify [post]
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

// ListUsers godoc
// @Summary List users with pagination
// @Description Returns a paginated list of users, with optional search and role filter.
// @Tags User
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)" minimum(1)
// @Param limit query int false "Items per page (default 10, max 100)" minimum(1) maximum(100)
// @Param search query string false "Search keyword (matches name, phone, email)"
// @Param role query string false "Role filter (default user)" Enums(admin, employee, user)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/user/list [get]
func (c *User) ListUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	search := ctx.Query("search")
	role := domain.Role(ctx.DefaultQuery("role", "user"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	users, total, err := c.userService.ListUsers(offset, limit, search, role)
	if err != nil {
		ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	data := make([]dto.User, len(users))
	for i, u := range users {
		data[i] = dto.User{
			Id:          u.ID.String(),
			PhoneNumber: u.PhoneNumber,
			Email:       u.Email,
			FullName:    u.FullName,
			Photo:       u.Photo,
			Role:        string(u.Role),
			Verified:    u.Verified,
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	ctx.JSON(200, dto.PaginatedResponse{
		Code:    200,
		Message: "success",
		Data:    data,
		Meta: dto.Meta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
