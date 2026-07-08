package service

import (
	"github.com/google/uuid"
	"github.com/omnlgy/jadwalin/internal/domain"
)

type TreatmentService struct {
	treatmentRepo domain.TreatmentRepository
}

func NewTreatmentService(treatmentRepo domain.TreatmentRepository) *TreatmentService {
	return &TreatmentService{treatmentRepo: treatmentRepo}
}

func (s *TreatmentService) Create(t *domain.Treatment) (*domain.Treatment, error) {
	return s.treatmentRepo.Create(t)
}

func (s *TreatmentService) GetByID(id uuid.UUID) (*domain.Treatment, error) {
	return s.treatmentRepo.GetByID(id)
}

func (s *TreatmentService) List(offset, limit int, search string) ([]domain.Treatment, int64, error) {
	return s.treatmentRepo.List(offset, limit, search)
}

func (s *TreatmentService) Update(t *domain.Treatment) error {
	return s.treatmentRepo.Update(t)
}

func (s *TreatmentService) Delete(id uuid.UUID) error {
	return s.treatmentRepo.Delete(id)
}
