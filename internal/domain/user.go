package domain

import "github.com/gofrs/uuid/v5"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
	RoleUser     Role = "user"
)

type User struct {
	ID          uuid.UUID
	PhoneNumber string
	Email       string
	Address     string
	FullName    string
	Photo       string
	Role        Role
}

type UserRepository interface {
	Create(user *User) (*User, error)
	GetByID(id uuid.UUID) (*User, error)
	GetByPhoneNumber(phoneNumber string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
}

type UserService interface {
	GetByID(id uuid.UUID) (*User, error)
	RegisterEmployee(user *User) (*User, error)
}
