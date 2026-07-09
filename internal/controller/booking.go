package controller

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omnlgy/jadwalin/internal/domain"
	"github.com/omnlgy/jadwalin/internal/dto"
)

type Booking struct {
	bookingService   domain.BookingService
	treatmentService domain.TreatmentService
}

func NewBookingController(bookingService domain.BookingService, treatmentService domain.TreatmentService) *Booking {
	return &Booking{
		bookingService:   bookingService,
		treatmentService: treatmentService,
	}
}

// GetAvailableSlots godoc
// @Summary Get available time slots for a treatment + staff + date
// @Description Returns all available time slots for a given treatment, staff, and date
// @Tags Booking
// @Accept json
// @Produce json
// @Param request body dto.GetAvailabelSlotRequest true "Get Available Slots Request"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.BadRequestResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/booking/available-slots [post]
func (c *Booking) GetAvailableSlots(ctx *gin.Context) {
	var body dto.GetAvailabelSlotRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		abortWithBadRequest(ctx, err)
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

	staffID, err := uuid.Parse(body.StaffID)
	if err != nil {
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: "invalid staff_id",
		})
		return
	}

	date, err := time.Parse("2006-01-02", body.Date)
	if err != nil {
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: "invalid date format, expected YYYY-MM-DD",
		})
		return
	}

	slots, err := c.bookingService.GetAvailableSlots(treatmentID, staffID, date)
	if err != nil {
		ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, dto.SuccessResponse{
		Code:    200,
		Message: "success",
		Data:    slots,
	})
}

// CreateBooking godoc
// @Summary Create a new booking
// @Description Creates a new booking (authenticated user)
// @Tags Booking
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.CreateBookingRequest true "Create Booking Request"
// @Success 201 {object} dto.CreatedResponse
// @Failure 400 {object} dto.BadRequestResponse
// @Failure 409 {object} dto.BadRequestResponse
// @Failure 500 {object} dto.InternalErrorResponse
// @Router /api/booking [post]
func (c *Booking) CreateBooking(ctx *gin.Context) {
	userIDStr, exists := ctx.Get("userId")
	if !exists {
		ctx.AbortWithStatusJSON(401, dto.UnauthorizedResponse{
			Code:    401,
			Message: "unauthorized",
		})
		return
	}
	clientID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
			Code:    500,
			Message: "failed to parse user id from token",
		})
		return
	}

	var body dto.CreateBookingRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		abortWithBadRequest(ctx, err)
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

	staffID, err := uuid.Parse(body.StaffID)
	if err != nil {
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: "invalid staff_id",
		})
		return
	}

	startTime, err := time.Parse(time.RFC3339, body.StartTime)
	if err != nil {
		ctx.AbortWithStatusJSON(400, dto.BadRequestResponse{
			Code:    400,
			Message: "invalid start_time format, expected RFC3339 (e.g. 2025-10-15T10:00:00Z)",
		})
		return
	}

	// Get treatment to calculate end_time
	treatment, err := c.treatmentService.GetByID(treatmentID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ctx.AbortWithStatusJSON(404, dto.BadRequestResponse{
				Code:    400,
				Message: "treatment not found",
			})
			return
		}
		ctx.AbortWithStatusJSON(500, dto.InternalErrorResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	endTime := startTime.Add(time.Duration(treatment.Duration) * time.Minute)

	booking := &domain.Booking{
		ClientID:    clientID,
		StaffID:     staffID,
		TreatmentID: treatmentID,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	result, err := c.bookingService.Create(booking)
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			ctx.AbortWithStatusJSON(409, dto.BadRequestResponse{
				Code:    409,
				Message: err.Error(),
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
		Message: "booking created",
		Data: dto.Booking{
			Id:          result.ID.String(),
			ClientID:    result.ClientID.String(),
			StaffID:     result.StaffID.String(),
			TreatmentID: result.TreatmentID.String(),
			StartTime:   result.StartTime.Format(time.RFC3339),
			EndTime:     result.EndTime.Format(time.RFC3339),
			Status:      result.Status,
		},
	})
}
