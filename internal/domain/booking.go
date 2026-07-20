package domain

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID          uuid.UUID  `json:"id"`
	ClientID    uuid.UUID  `json:"client_id"`
	Client      *User      `json:"client"`
	TreatmentID uuid.UUID  `json:"treatment_id"`
	Treatment   *Treatment `json:"treatment"`
	StaffID     uuid.UUID  `json:"staff_id"`
	Staff       *User      `json:"staff"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     time.Time  `json:"end_time"`
	Status      string     `json:"status"`
}

type AvailableSlot struct {
	StartTime time.Time
	EndTime   time.Time
}

type BookingQuery struct {
	Status        string
	TreatmentName string
}

type BookingRepository interface {
	Create(booking *Booking) (*Booking, error)
	GetByID(id uuid.UUID) (*Booking, error)
	GetByUserID(userID uuid.UUID, params BookingQuery) ([]Booking, error)
	GetByTreatmentID(treatmentID uuid.UUID) ([]Booking, error)
	GetByStaffID(staffID uuid.UUID) ([]Booking, error)
	GetByDate(date time.Time) ([]Booking, error)
	Update(booking *Booking) (*Booking, error)
	Delete(id uuid.UUID) error
	GetBookingsForStaffAndDate(staffID uuid.UUID, date time.Time) ([]Booking, error)
	GetBookingInRange(startTime, endTime time.Time) ([]Booking, error)
}

type BookingService interface {
	Create(booking *Booking) (*Booking, error)
	GetByID(id uuid.UUID) (*Booking, error)
	GetByUserID(userID uuid.UUID, params BookingQuery) ([]Booking, error)
	GetByTreatmentID(treatmentID uuid.UUID) ([]Booking, error)
	GetByStaffID(staffID uuid.UUID) ([]Booking, error)
	GetByDate(date time.Time) ([]Booking, error)
	Update(booking *Booking) (*Booking, error)
	Delete(id uuid.UUID) error
	GetBookingsForStaffAndDate(staffID uuid.UUID, date time.Time) ([]Booking, error)
	GetAvailableSlots(treatmentID, staffID uuid.UUID, date time.Time) ([]AvailableSlot, error)
}
