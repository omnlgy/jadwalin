package domain

import "github.com/gofrs/uuid/v5"

type User struct {
	ID          uuid.UUID
	PhoneNumber string
	Address     string
	FullName    string
	Photo       string
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByPhoneNumber(phoneNumber string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
}
