package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/omnlgy/jadwalin/internal/domain"
	"github.com/omnlgy/jadwalin/internal/models"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *User {
	return &User{db: db}
}

func (r *User) Create(user *domain.User) (*domain.User, error) {
	m := &models.User{
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Address:     user.Address,
		FullName:    user.FullName,
		Photo:       user.Photo,
		Role:        string(user.Role),
	}
	if err := r.db.Create(m).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || isPgUniqueViolation(err) {
			return nil, fmt.Errorf("repo: create user: %w", domain.ErrConflict)
		}
		return nil, fmt.Errorf("repo: create user: %w", err)
	}
	user.ID = m.ID
	return user, nil
}

func (r *User) GetByID(id uuid.UUID) (*domain.User, error) {
	var m models.User
	err := r.db.First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("repo: get user %s: %w", id, domain.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("repo: get user %s: %w", id, err)
	}
	return toDomainUser(&m), nil
}

func (r *User) GetByPhoneNumber(phoneNumber string) (*domain.User, error) {
	var m models.User
	err := r.db.Where("phone_number = ?", phoneNumber).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("repo: get user by phone %s: %w", phoneNumber, domain.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("repo: get user by phone %s: %w", phoneNumber, err)
	}
	return toDomainUser(&m), nil
}

func (r *User) Update(user *domain.User) error {
	m := &models.User{
		ID:          user.ID,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Address:     user.Address,
		FullName:    user.FullName,
		Photo:       user.Photo,
		Role:        string(user.Role),
	}
	result := r.db.Model(&models.User{}).Where("id = ?", user.ID).Updates(m)
	if result.Error != nil {
		return fmt.Errorf("repo: update user %s: %w", user.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("repo: update user %s: %w", user.ID, domain.ErrNotFound)
	}
	return nil
}

func (r *User) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("repo: delete user %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("repo: delete user %s: %w", id, domain.ErrNotFound)
	}
	return nil
}

func (r *User) List(offset, limit int, search string, role domain.Role) ([]domain.User, int64, error) {
	var total int64
	query := r.db.Model(&models.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("full_name ILIKE ? OR phone_number ILIKE ? OR email ILIKE ?", like, like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("repo: count users: %w", err)
	}

	if total == 0 {
		return []domain.User{}, 0, nil
	}

	var ms []models.User
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&ms).Error; err != nil {
		return nil, 0, fmt.Errorf("repo: list users: %w", err)
	}

	users := make([]domain.User, len(ms))
	for i, m := range ms {
		users[i] = *toDomainUser(&m)
	}
	return users, total, nil
}

func isPgUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func toDomainUser(m *models.User) *domain.User {
	return &domain.User{
		ID:          m.ID,
		PhoneNumber: m.PhoneNumber,
		Email:       m.Email,
		Address:     m.Address,
		FullName:    m.FullName,
		Photo:       m.Photo,
		Role:        domain.Role(m.Role),
		Verified:    m.Verified,
	}
}
