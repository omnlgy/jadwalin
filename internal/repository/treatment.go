package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/omnlgy/jadwalin/internal/domain"
	"github.com/omnlgy/jadwalin/internal/models"
	"gorm.io/gorm"
)

type Treatment struct {
	db *gorm.DB
}

func NewTreatmentRepository(db *gorm.DB) *Treatment {
	return &Treatment{db: db}
}

func (r *Treatment) Create(t *domain.Treatment) (*domain.Treatment, error) {
	m := &models.Treatment{
		Name:        t.Name,
		Description: t.Description,
		Duration:    t.Duration,
		Price:       t.Price,
	}
	if err := r.db.Create(m).Error; err != nil {
		return nil, fmt.Errorf("repo: create treatment: %w", err)
	}
	t.ID = m.ID
	return t, nil
}

func (r *Treatment) GetByID(id uuid.UUID) (*domain.Treatment, error) {
	var m models.Treatment
	err := r.db.First(&m, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("repo: get treatment %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("repo: get treatment %s: %w", id, err)
	}
	return toDomainTreatment(&m), nil
}

func (r *Treatment) List(offset, limit int, search string) ([]domain.Treatment, int64, error) {
	var total int64
	query := r.db.Model(&models.Treatment{})

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ?", like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("repo: count treatments: %w", err)
	}

	if total == 0 {
		return []domain.Treatment{}, 0, nil
	}

	var ms []models.Treatment
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&ms).Error; err != nil {
		return nil, 0, fmt.Errorf("repo: list treatments: %w", err)
	}

	treatments := make([]domain.Treatment, len(ms))
	for i, m := range ms {
		treatments[i] = *toDomainTreatment(&m)
	}
	return treatments, total, nil
}

func (r *Treatment) Update(t *domain.Treatment) error {
	result := r.db.Model(&models.Treatment{}).Where("id = ?", t.ID).Updates(map[string]interface{}{
		"name":        t.Name,
		"description": t.Description,
		"duration":    t.Duration,
		"price":       t.Price,
	})
	if result.Error != nil {
		return fmt.Errorf("repo: update treatment %s: %w", t.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("repo: update treatment %s: %w", t.ID, domain.ErrNotFound)
	}
	return nil
}

func (r *Treatment) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Treatment{}, id)
	if result.Error != nil {
		return fmt.Errorf("repo: delete treatment %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("repo: delete treatment %s: %w", id, domain.ErrNotFound)
	}
	return nil
}

func toDomainTreatment(m *models.Treatment) *domain.Treatment {
	return &domain.Treatment{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Duration:    m.Duration,
		Price:       m.Price,
	}
}
