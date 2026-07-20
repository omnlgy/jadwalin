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

type StaffSkill struct {
	staffSkillService domain.StaffSkillService
}

func NewStaffSkillController(staffSkillService domain.StaffSkillService) *StaffSkill {
	return &StaffSkill{staffSkillService: staffSkillService}
}

// AssignSkill godoc
// @Summary Assign a treatment skill to a staff
// @Description Assigns a treatment to a staff member (admin only)
// @Tags StaffSkill
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.AssignSkillRequest true "Assign Skill Request"
// @Success 201 {object} dto.CreatedResponse
// @Failure 400 {object} dto.BadRequestResponse
// @Failure 409 {object} dto.BadRequestResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/staff-skills [post]
func (c *StaffSkill) AssignSkill(ctx *gin.Context) {
	var body dto.AssignSkillRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		abortWithBadRequest(ctx, err)
		return
	}

	userID, err := uuid.Parse(body.UserID)
	if err != nil {
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: "invalid user_id",
		})
		return
	}

	treatmentID, err := uuid.Parse(body.TreatmentID)
	if err != nil {
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: "invalid treatment_id",
		})
		return
	}

	result, err := c.staffSkillService.Assign(userID, treatmentID)
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			ctx.AbortWithStatusJSON(409, dto.BadRequestResponse{
				Code:    409,
				Message: "staff already has this skill assigned",
			})
			return
		}
		ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(201, dto.CreatedResponse{
		Code:    201,
		Message: "Skill assigned successfully",
		Data:    toStaffSkillDTO(result),
	})
}

// UnassignSkill godoc
// @Summary Remove a skill from a staff
// @Description Unassigns a treatment from a staff member (admin only)
// @Tags StaffSkill
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "StaffSkill ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.BadRequestResponse
// @Failure 404 {object} dto.BadRequestResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/staff-skills/{id} [delete]
func (c *StaffSkill) UnassignSkill(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: "invalid staff skill id",
		})
		return
	}

	if err := c.staffSkillService.Unassign(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ctx.AbortWithStatusJSON(404, dto.BadRequestResponse{
				Code:    400,
				Message: "staff skill not found",
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
		Message: "skill unassigned",
	})
}

// GetStaffSkill godoc
// @Summary Get a staff skill by ID
// @Description Returns a single staff-skill assignment with joined user and treatment
// @Tags StaffSkill
// @Produce json
// @Param id path string true "StaffSkill ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.BadRequestResponse
// @Failure 404 {object} dto.BadRequestResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/staff-skills/{id} [get]
func (c *StaffSkill) GetStaffSkill(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: "invalid staff skill id",
		})
		return
	}

	ss, err := c.staffSkillService.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ctx.AbortWithStatusJSON(404, dto.BadRequestResponse{
				Code:    400,
				Message: "staff skill not found",
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
		Data:    toStaffSkillDTO(ss),
	})
}

// ListByStaff godoc
// @Summary List skills for a staff member
// @Description Returns all treatments a staff member can perform
// @Tags StaffSkill
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.BadRequestResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/staff-skills/staff/{userId} [get]
func (c *StaffSkill) ListByStaff(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("userId"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: "invalid user id",
		})
		return
	}

	skills, err := c.staffSkillService.ListByStaff(userID)
	if err != nil {
		ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	data := make([]dto.StaffSkill, len(skills))
	for i, s := range skills {
		data[i] = toStaffSkillDTO(&s)
	}

	ctx.JSON(200, dto.SuccessResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// ListByTreatment godoc
// @Summary List staff who can perform a treatment
// @Description Returns all staff members skilled at a given treatment
// @Tags StaffSkill
// @Produce json
// @Param treatmentId path string true "Treatment ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.BadRequestResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/staff-skills/treatment/{treatmentId} [get]
func (c *StaffSkill) ListByTreatment(ctx *gin.Context) {
	treatmentID, err := uuid.Parse(ctx.Param("treatmentId"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: "invalid treatment id",
		})
		return
	}

	skills, err := c.staffSkillService.ListByTreatment(treatmentID)
	if err != nil {
		ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	data := make([]dto.StaffSkill, len(skills))
	for i, s := range skills {
		data[i] = toStaffSkillDTO(&s)
	}

	ctx.JSON(200, dto.SuccessResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// ListAll godoc
// @Summary List all staff-skill assignments with pagination
// @Description Returns a paginated list of all staff-skill assignments, with optional search by staff name
// @Tags StaffSkill
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)" minimum(1)
// @Param limit query int false "Items per page (default 10, max 100)" minimum(1) maximum(100)
// @Param search query string false "Search keyword (matches staff name or phone)"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/staff-skills/list [get]
func (c *StaffSkill) ListAll(ctx *gin.Context) {
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

	skills, total, err := c.staffSkillService.ListAll(offset, limit, search)
	if err != nil {
		ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	data := make([]dto.StaffSkill, len(skills))
	for i, s := range skills {
		data[i] = toStaffSkillDTO(&s)
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

func toStaffSkillDTO(s *domain.StaffSkill) dto.StaffSkill {
	d := dto.StaffSkill{
		Id:          s.ID.String(),
		UserID:      s.UserID.String(),
		TreatmentID: s.TreatmentID.String(),
	}
	if s.User != nil {
		d.UserFullName = s.User.FullName
		d.UserPhoneNumber = s.User.PhoneNumber
	}
	if s.Treatment != nil {
		d.TreatmentName = s.Treatment.Name
		d.TreatmentPrice = s.Treatment.Price
		d.TreatmentDuration = s.Treatment.Duration
	}
	return d
}
