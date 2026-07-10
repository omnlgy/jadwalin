package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Booking struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	ClientID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	StaffID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	TreatmentID uuid.UUID      `gorm:"type:uuid;not null;index"`
	StartTime   time.Time      `gorm:"not null"`
	EndTime     time.Time      `gorm:"not null"`
	Status      string         `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending','confirmed','completed','cancelled')"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Associations
	Client    User      `gorm:"foreignKey:ClientID;references:ID"`
	Staff     User      `gorm:"foreignKey:StaffID;references:ID"`
	Treatment Treatment `gorm:"foreignKey:TreatmentID;references:ID"`
}

// BeforeCreate will set a UUIDV7 for the ID field
func (b *Booking) BeforeCreate(tx *gorm.DB) (err error) {
	newUUID, err := uuid.NewV7()
	if err != nil {
		return err
	}
	b.ID = newUUID
	return nil
}
