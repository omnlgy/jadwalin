package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StaffSkill struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	User        User      `gorm:"foreignKey:UserID"`
	TreatmentID uuid.UUID `gorm:"type:uuid;not null;index"`
	Treatment   Treatment `gorm:"foreignKey:TreatmentID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// BeforeCreate will set a UUIDV7 for the ID field
func (ss *StaffSkill) BeforeCreate(tx *gorm.DB) (err error) {
	newUUID, err := uuid.NewV7()
	if err != nil {
		return err
	}
	ss.ID = newUUID
	return nil
}
