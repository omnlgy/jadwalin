package service

import (
	"github.com/google/uuid"
	"github.com/omnlgy/jadwalin/internal/domain"
)

type StaffSkillService struct {
	repo domain.StaffSkillRepository
}

func NewStaffSkillService(repo domain.StaffSkillRepository) *StaffSkillService {
	return &StaffSkillService{repo: repo}
}

func (s *StaffSkillService) Assign(userID, treatmentID uuid.UUID) (*domain.StaffSkill, error) {
	return s.repo.Assign(userID, treatmentID)
}

func (s *StaffSkillService) Unassign(id uuid.UUID) error {
	return s.repo.Unassign(id)
}

func (s *StaffSkillService) GetByID(id uuid.UUID) (*domain.StaffSkill, error) {
	return s.repo.GetByID(id)
}

func (s *StaffSkillService) ListByStaff(userID uuid.UUID) ([]domain.StaffSkill, error) {
	return s.repo.ListByStaff(userID)
}

func (s *StaffSkillService) ListByTreatment(treatmentID uuid.UUID) ([]domain.StaffSkill, error) {
	return s.repo.ListByTreatment(treatmentID)
}

func (s *StaffSkillService) ListAll(offset, limit int, search string) ([]domain.StaffSkill, int64, error) {
	return s.repo.ListAll(offset, limit, search)
}
