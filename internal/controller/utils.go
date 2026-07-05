package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/omnlgy/jadwalin/internal/dto"
	"github.com/omnlgy/jadwalin/utils"
)

func abortWithBadRequest(ctx *gin.Context, err error) {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		fieldErrs := make([]dto.FieldError, len(validationErrs))

		for i, fe := range validationErrs {
			fieldErrs[i] = dto.FieldError{
				Field:   utils.PascalToSnake(fe.Field()),
				Message: fe.Tag(),
			}
		}
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: "Validation failed",
			Errors:  fieldErrs,
		})
	} else {
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: err.Error(),
		})
	}

}
