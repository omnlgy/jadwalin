package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Treatment struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name        string         `gorm:"not null"`
	Description string
	Duration    int            `gorm:"not null"`
	Price       float64        `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate will set a UUIDV7 for the ID field
func (t *Treatment) BeforeCreate(tx *gorm.DB) (err error) {
	newUUID, err := uuid.NewV7()
	if err != nil {
		return err
	}
	t.ID = newUUID
	return nil
}
