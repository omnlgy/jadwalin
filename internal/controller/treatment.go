package controller

import (
	"errors"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omnlgy/jadwalin/internal/domain"
	"github.com/omnlgy/jadwalin/internal/dto"
)

type Treatment struct {
	treatmentService domain.TreatmentService
}

func NewTreatmentController(treatmentService domain.TreatmentService) *Treatment {
	return &Treatment{treatmentService: treatmentService}
}

// CreateTreatment godoc
// @Summary Create a new treatment
// @Description Creates a new treatment (admin only)
// @Tags Treatment
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.CreateTreatmentRequest true "Create Treatment Request"
// @Success 201 {object} dto.CreatedResponse
// @Failure 400 {object} dto.BadRequestResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/treatment [post]
func (c *Treatment) CreateTreatment(ctx *gin.Context) {
	var body dto.CreateTreatmentRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		abortWithBadRequest(ctx, err)
		return
	}

	treatment := &domain.Treatment{
		Name:        body.Name,
		Description: body.Description,
		Duration:    body.Duration,
		Price:       body.Price,
	}

	result, err := c.treatmentService.Create(treatment)
	if err != nil {
		ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(201, dto.CreatedResponse{
		Code:    201,
		Message: "Treatment created successfully",
		Data: dto.Treatment{
			Id:          result.ID.String(),
			Name:        result.Name,
			Description: result.Description,
			Duration:    result.Duration,
			Price:       result.Price,
		},
	})
}

// GetTreatment godoc
// @Summary Get a treatment by ID
// @Description Returns a single treatment
// @Tags Treatment
// @Produce json
// @Param id path string true "Treatment ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 404 {object} dto.BadRequestResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/treatment/{id} [get]
func (c *Treatment) GetTreatment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: "invalid treatment id",
		})
		return
	}

	treatment, err := c.treatmentService.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ctx.AbortWithStatusJSON(404, dto.BadRequestResponse{
				Code:    400,
				Message: "treatment not found",
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
		Message: "success",
		Data: dto.Treatment{
			Id:          treatment.ID.String(),
			Name:        treatment.Name,
			Description: treatment.Description,
			Duration:    treatment.Duration,
			Price:       treatment.Price,
		},
	})
}

// ListTreatments godoc
// @Summary List treatments with pagination
// @Description Returns a paginated list of treatments with optional search
// @Tags Treatment
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)" minimum(1)
// @Param limit query int false "Items per page (default 10, max 100)" minimum(1) maximum(100)
// @Param search query string false "Search keyword (matches name)"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/treatment/list [get]
func (c *Treatment) ListTreatments(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	search := ctx.Query("search")

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

	treatments, total, err := c.treatmentService.List(offset, limit, search)
	if err != nil {
		ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	data := make([]dto.Treatment, len(treatments))
	for i, t := range treatments {
		data[i] = dto.Treatment{
			Id:          t.ID.String(),
			Name:        t.Name,
			Description: t.Description,
			Duration:    t.Duration,
			Price:       t.Price,
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

// UpdateTreatment godoc
// @Summary Update a treatment
// @Description Updates an existing treatment
// @Tags Treatment
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Treatment ID"
// @Param request body dto.UpdateTreatmentRequest true "Update Treatment Request"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.BadRequestResponse
// @Failure 404 {object} dto.BadRequestResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/treatment/{id} [put]
func (c *Treatment) UpdateTreatment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: "invalid treatment id",
		})
		return
	}

	current, err := c.treatmentService.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ctx.AbortWithStatusJSON(404, dto.BadRequestResponse{
				Code:    400,
				Message: "treatment not found",
			})
		} else {
			ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}
		return
	}

	var body dto.UpdateTreatmentRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		abortWithBadRequest(ctx, err)
		return
	}

	if body.Name != "" {
		current.Name = body.Name
	}
	if body.Description != "" {
		current.Description = body.Description
	}
	if body.Duration != 0 {
		current.Duration = body.Duration
	}
	if body.Price != 0 {
		current.Price = body.Price
	}

	if err := c.treatmentService.Update(current); err != nil {
		ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, dto.SuccessResponse{
		Code:    200,
		Message: "treatment updated",
	})
}

// DeleteTreatment godoc
// @Summary Delete a treatment
// @Description Deletes a treatment by ID (soft delete)
// @Tags Treatment
// @Security ApiKeyAuth
// @Param id path string true "Treatment ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.BadRequestResponse
// @Failure 404 {object} dto.BadRequestResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/treatment/{id} [delete]
func (c *Treatment) DeleteTreatment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: "invalid treatment id",
		})
		return
	}

	if err := c.treatmentService.Delete(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ctx.AbortWithStatusJSON(404, dto.BadRequestResponse{
				Code:    400,
				Message: "treatment not found",
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
		Message: "treatment deleted",
	})
}
