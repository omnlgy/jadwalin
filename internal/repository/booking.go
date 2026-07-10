package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/omnlgy/jadwalin/internal/domain"
	"github.com/omnlgy/jadwalin/internal/models"
	"gorm.io/gorm"
)

type Booking struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *Booking {
	return &Booking{db: db}
}

func (r *Booking) Create(booking *domain.Booking) (*domain.Booking, error) {
	m := &models.Booking{
		ClientID:    booking.ClientID,
		StaffID:     booking.StaffID,
		TreatmentID: booking.TreatmentID,
		StartTime:   booking.StartTime,
		EndTime:     booking.EndTime,
		Status:      booking.Status,
	}
	if err := r.db.Create(m).Error; err != nil {
		return nil, fmt.Errorf("repo: create booking: %w", err)
	}
	created, err := r.GetByID(m.ID)
	if err != nil {
		return nil, fmt.Errorf("repo: create booking: failed to reload: %w", err)
	}
	return created, nil
}

func (r *Booking) GetByID(id uuid.UUID) (*domain.Booking, error) {
	var m models.Booking
	err := r.db.First(&m, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("repo: get booking %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("repo: get booking %s: %w", id, err)
	}
	// Load associations manually to avoid GORM dual-FK-to-same-table bug
	var client models.User
	r.db.Unscoped().First(&client, m.ClientID)
	m.Client = client
	var staff models.User
	r.db.Unscoped().First(&staff, m.StaffID)
	m.Staff = staff
	var treatment models.Treatment
	r.db.First(&treatment, m.TreatmentID)
	m.Treatment = treatment
	return toDomainBooking(&m), nil
}

func (r *Booking) GetByUserID(userID uuid.UUID) ([]domain.Booking, error) {
	var ms []models.Booking
	err := r.db.Where("client_id = ?", userID).
		Preload("Client").Preload("Staff").Preload("Treatment").
		Order("start_time DESC").Find(&ms).Error
	if err != nil {
		return nil, fmt.Errorf("repo: get bookings by user %s: %w", userID, err)
	}
	return sliceToDomainBooking(ms), nil
}

func (r *Booking) GetByTreatmentID(treatmentID uuid.UUID) ([]domain.Booking, error) {
	var ms []models.Booking
	err := r.db.Where("treatment_id = ?", treatmentID).
		Preload("Client").Preload("Staff").Preload("Treatment").
		Order("start_time DESC").Find(&ms).Error
	if err != nil {
		return nil, fmt.Errorf("repo: get bookings by treatment %s: %w", treatmentID, err)
	}
	return sliceToDomainBooking(ms), nil
}

func (r *Booking) GetByStaffID(staffID uuid.UUID) ([]domain.Booking, error) {
	var ms []models.Booking
	err := r.db.Where("staff_id = ?", staffID).
		Preload("Client").Preload("Staff").Preload("Treatment").
		Order("start_time DESC").Find(&ms).Error
	if err != nil {
		return nil, fmt.Errorf("repo: get bookings by staff %s: %w", staffID, err)
	}
	return sliceToDomainBooking(ms), nil
}

func (r *Booking) GetByDate(date time.Time) ([]domain.Booking, error) {
	var ms []models.Booking
	startOfDay := date.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)
	err := r.db.Where("start_time >= ? AND start_time < ?", startOfDay, endOfDay).
		Preload("Client").Preload("Staff").Preload("Treatment").
		Order("start_time ASC").Find(&ms).Error
	if err != nil {
		return nil, fmt.Errorf("repo: get bookings by date: %w", err)
	}
	return sliceToDomainBooking(ms), nil
}

func (r *Booking) Update(booking *domain.Booking) (*domain.Booking, error) {
	m := &models.Booking{
		ID:          booking.ID,
		ClientID:    booking.ClientID,
		StaffID:     booking.StaffID,
		TreatmentID: booking.TreatmentID,
		StartTime:   booking.StartTime,
		EndTime:     booking.EndTime,
		Status:      booking.Status,
	}
	result := r.db.Model(&models.Booking{}).Where("id = ?", booking.ID).Updates(m)
	if result.Error != nil {
		return nil, fmt.Errorf("repo: update booking %s: %w", booking.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("repo: update booking %s: %w", booking.ID, domain.ErrNotFound)
	}
	return r.GetByID(booking.ID)
}

func (r *Booking) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Booking{}, id)
	if result.Error != nil {
		return fmt.Errorf("repo: delete booking %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("repo: delete booking %s: %w", id, domain.ErrNotFound)
	}
	return nil
}

func (r *Booking) GetBookingsForStaffAndDate(staffID uuid.UUID, date time.Time) ([]domain.Booking, error) {
	var bookings []models.Booking
	startOfDay := date.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)
	err := r.db.
		Where("staff_id = ? AND start_time >= ? AND start_time < ? AND status != ?", staffID, startOfDay, endOfDay, "cancelled").
		Order("start_time ASC").
		Find(&bookings).Error

	if err != nil {
		return nil, fmt.Errorf("repo: get bookings for staff and date: %w", err)
	}

	bookingsDomain := make([]domain.Booking, len(bookings))
	for i, booking := range bookings {
		bookingsDomain[i] = *toDomainBooking(&booking)
	}
	return bookingsDomain, nil
}

func (r *Booking) GetBookingInRange(startTime, endTime time.Time) ([]domain.Booking, error) {
	var bookings []models.Booking
	err := r.db.
		Where("start_time >= ? AND start_time < ? AND status != ?", startTime, endTime, "cancelled").
		Order("start_time ASC").
		Find(&bookings).Error

	if err != nil {
		return nil, fmt.Errorf("repo: get bookings in range: %w", err)
	}

	bookingsDomain := make([]domain.Booking, len(bookings))
	for i, booking := range bookings {
		bookingsDomain[i] = *toDomainBooking(&booking)
	}
	return bookingsDomain, nil
}

func toDomainBooking(m *models.Booking) *domain.Booking {
	db := &domain.Booking{
		ID:          m.ID,
		ClientID:    m.ClientID,
		StaffID:     m.StaffID,
		TreatmentID: m.TreatmentID,
		StartTime:   m.StartTime,
		EndTime:     m.EndTime,
		Status:      m.Status,
	}

	if m.ClientID != uuid.Nil {
		db.Client = toDomainUser(&m.Client)
	}

	if m.StaffID != uuid.Nil {
		db.Staff = toDomainUser(&m.Staff)
	}

	if m.TreatmentID != uuid.Nil {
		db.Treatment = toDomainTreatment(&m.Treatment)
	}

	return db
}

func sliceToDomainBooking(ms []models.Booking) []domain.Booking {
	res := make([]domain.Booking, len(ms))
	for i, m := range ms {
		res[i] = *toDomainBooking(&m)
	}
	return res
}
