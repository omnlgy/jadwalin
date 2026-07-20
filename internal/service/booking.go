package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/omnlgy/jadwalin/internal/domain"
)

type BookingService struct {
	bookingRepo    domain.BookingRepository
	userRepo       domain.UserRepository
	treatmentRepo  domain.TreatmentRepository
	staffSkillRepo domain.StaffSkillRepository
}

func NewBookingService(bookingRepo domain.BookingRepository, userRepo domain.UserRepository, treatmentRepo domain.TreatmentRepository, staffSkillRepo domain.StaffSkillRepository) *BookingService {
	return &BookingService{
		bookingRepo:    bookingRepo,
		userRepo:       userRepo,
		treatmentRepo:  treatmentRepo,
		staffSkillRepo: staffSkillRepo,
	}
}

func (s *BookingService) Create(booking *domain.Booking) (*domain.Booking, error) {
	staffSkill, err := s.staffSkillRepo.GetByStaffAndTreatment(booking.StaffID, booking.TreatmentID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("service: staff is not skilled for this treatment: %w", domain.ErrConflict)
		}
		return nil, fmt.Errorf("service: check staff skill: %w", err)
	}
	if staffSkill == nil {
		return nil, fmt.Errorf("service: staff is not skilled for this treatment: %w", domain.ErrConflict)
	}

	bookings, err := s.bookingRepo.GetBookingInRange(booking.StartTime, booking.EndTime)
	if err != nil {
		return nil, fmt.Errorf("service: check availability: %w", err)
	}

	if len(bookings) > 0 {
		return nil, fmt.Errorf("service: slot not available: %w", domain.ErrConflict)
	}

	booking.Status = "pending"

	return s.bookingRepo.Create(booking)
}

func (s *BookingService) GetByID(id uuid.UUID) (*domain.Booking, error) {
	return s.bookingRepo.GetByID(id)
}

func (s *BookingService) GetByUserID(userID uuid.UUID, params domain.BookingQuery) ([]domain.Booking, error) {
	return s.bookingRepo.GetByUserID(userID, params)
}

func (s *BookingService) GetByTreatmentID(treatmentID uuid.UUID) ([]domain.Booking, error) {
	return s.bookingRepo.GetByTreatmentID(treatmentID)
}

func (s *BookingService) GetByStaffID(staffID uuid.UUID) ([]domain.Booking, error) {
	return s.bookingRepo.GetByStaffID(staffID)
}

func (s *BookingService) GetByDate(date time.Time) ([]domain.Booking, error) {
	return s.bookingRepo.GetByDate(date)
}

func (s *BookingService) Update(booking *domain.Booking) (*domain.Booking, error) {
	return s.bookingRepo.Update(booking)
}

func (s *BookingService) Delete(id uuid.UUID) error {
	return s.bookingRepo.Delete(id)
}

func (s *BookingService) GetBookingsForStaffAndDate(staffID uuid.UUID, date time.Time) ([]domain.Booking, error) {
	return s.bookingRepo.GetBookingsForStaffAndDate(staffID, date)
}

func (s *BookingService) GetAvailableSlots(treatmentID, staffID uuid.UUID, date time.Time) ([]domain.AvailableSlot, error) {
	treatment, err := s.treatmentRepo.GetByID(treatmentID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get treatment duration: %w", err)
	}
	treatmentDuration := time.Duration(treatment.Duration) * time.Minute

	storeOpen := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, date.Location())
	storeClose := time.Date(date.Year(), date.Month(), date.Day(), 17, 0, 0, 0, date.Location())

	existingBookings, err := s.bookingRepo.GetBookingsForStaffAndDate(staffID, date)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get existing bookings: %w", err)
	}

	var availableSlots []domain.AvailableSlot
	slotStartTime := storeOpen

	now := time.Now().In(date.Location())
	if slotStartTime.Before(now) {
		slotStartTime = now
	}

	for !slotStartTime.Add(treatmentDuration).After(storeClose) {
		slotEndTime := slotStartTime.Add(treatmentDuration)
		isConflict := false
		nextTime := slotStartTime

		for _, booking := range existingBookings {
			// booking start time is after or equal slot end time
			if !booking.StartTime.Before(slotEndTime) {
				break
			}

			// Strict overlap check
			if booking.EndTime.After(slotStartTime) {
				isConflict = true
				if booking.EndTime.After(nextTime) {
					nextTime = booking.EndTime
				}
			}
		}

		if !isConflict {
			availableSlots = append(availableSlots, domain.AvailableSlot{
				StartTime: slotStartTime,
				EndTime:   slotEndTime,
			})
			slotStartTime = slotEndTime // Advance by duration for back-to-back availability
		} else {
			slotStartTime = nextTime // Jump directly past the conflicting booking
		}
	}

	return availableSlots, nil
}
