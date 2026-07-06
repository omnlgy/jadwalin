package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	PhoneNumber string    `gorm:"unique;not null"`
	Email       string    `gorm:"unique;not null"`
	Address     string
	FullName    string
	Photo       string
	Role        string
	Verified    bool `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate will set a UUIDV7 for the ID field
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	newUUID, err := uuid.NewV7()
	if err != nil {
		return err
	}
	user.ID = newUUID
	return nil
}
