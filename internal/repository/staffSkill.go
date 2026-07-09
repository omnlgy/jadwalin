package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/omnlgy/jadwalin/internal/domain"
	"github.com/omnlgy/jadwalin/internal/models"
	"gorm.io/gorm"
)

type StaffSkill struct {
	db *gorm.DB
}

func NewStaffSkillRepository(db *gorm.DB) *StaffSkill {
	return &StaffSkill{db: db}
}

func (r *StaffSkill) Assign(userID, treatmentID uuid.UUID) (*domain.StaffSkill, error) {
	// Check if a non-deleted record already exists
	var existing models.StaffSkill
	err := r.db.Where("user_id = ? AND treatment_id = ?", userID, treatmentID).First(&existing).Error
	if err == nil {
		return nil, fmt.Errorf("repo: assign skill: %w", domain.ErrConflict)
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("repo: assign skill: %w", err)
	}

	m := &models.StaffSkill{
		UserID:      userID,
		TreatmentID: treatmentID,
	}
	if err := r.db.Create(m).Error; err != nil {
		return nil, fmt.Errorf("repo: assign skill: %w", err)
	}

	// Reload with Preload to get joined user and treatment
	var loaded models.StaffSkill
	if err := r.db.Preload("User").Preload("Treatment").First(&loaded, m.ID).Error; err != nil {
		return nil, fmt.Errorf("repo: assign skill: %w", err)
	}
	return toDomainStaffSkill(&loaded), nil
}

func (r *StaffSkill) Unassign(id uuid.UUID) error {
	result := r.db.Delete(&models.StaffSkill{}, id)
	if result.Error != nil {
		return fmt.Errorf("repo: unassign skill %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("repo: unassign skill %s: %w", id, domain.ErrNotFound)
	}
	return nil
}

func (r *StaffSkill) GetByID(id uuid.UUID) (*domain.StaffSkill, error) {
	var m models.StaffSkill
	err := r.db.Preload("User").Preload("Treatment").First(&m, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("repo: get staff skill %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("repo: get staff skill %s: %w", id, err)
	}
	return toDomainStaffSkill(&m), nil
}

func (r *StaffSkill) ListByStaff(userID uuid.UUID) ([]domain.StaffSkill, error) {
	var ms []models.StaffSkill
	err := r.db.Where("user_id = ?", userID).
		Preload("User").Preload("Treatment").
		Find(&ms).Error
	if err != nil {
		return nil, fmt.Errorf("repo: list skills by staff %s: %w", userID, err)
	}
	return sliceToDomainStaffSkill(ms), nil
}

func (r *StaffSkill) ListByTreatment(treatmentID uuid.UUID) ([]domain.StaffSkill, error) {
	var ms []models.StaffSkill
	err := r.db.Where("treatment_id = ?", treatmentID).
		Preload("User").Preload("Treatment").
		Find(&ms).Error
	if err != nil {
		return nil, fmt.Errorf("repo: list skills by treatment %s: %w", treatmentID, err)
	}
	return sliceToDomainStaffSkill(ms), nil
}

func (r *StaffSkill) ListAll(offset, limit int, search string) ([]domain.StaffSkill, int64, error) {
	var total int64
	query := r.db.Model(&models.StaffSkill{})

	if search != "" {
		like := "%" + search + "%"
		query = query.Joins("JOIN users ON users.id = staff_skills.user_id").
			Where("users.full_name ILIKE ? OR users.phone_number ILIKE ?", like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("repo: count staff skills: %w", err)
	}
	if total == 0 {
		return []domain.StaffSkill{}, 0, nil
	}

	var ms []models.StaffSkill
	if err := r.db.Preload("User").Preload("Treatment").
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&ms).Error; err != nil {
		return nil, 0, fmt.Errorf("repo: list staff skills: %w", err)
	}

	return sliceToDomainStaffSkill(ms), total, nil
}

func toDomainStaffSkill(m *models.StaffSkill) *domain.StaffSkill {
	ds := &domain.StaffSkill{
		ID:          m.ID,
		UserID:      m.UserID,
		TreatmentID: m.TreatmentID,
	}
	if m.User.ID != uuid.Nil {
		ds.User = toDomain(&m.User)
	}
	if m.Treatment.ID != uuid.Nil {
		ds.Treatment = toDomainTreatment(&m.Treatment)
	}
	return ds
}

func sliceToDomainStaffSkill(ms []models.StaffSkill) []domain.StaffSkill {
	res := make([]domain.StaffSkill, len(ms))
	for i, m := range ms {
		res[i] = *toDomainStaffSkill(&m)
	}
	return res
}
