package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PhoneNumber string         `gorm:"unique;not null" json:"phone_number"`
	Address     string         `json:"address"`
	FullName    string         `json:"full_name"`
	Photo       string         `json:"photo"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
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
