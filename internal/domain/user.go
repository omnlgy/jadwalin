package domain

import "github.com/google/uuid"

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
	Verified    bool
}

type UserRepository interface {
	Create(user *User) (*User, error)
	GetByID(id uuid.UUID) (*User, error)
	GetByPhoneNumber(phoneNumber string) (*User, error)
	List(offset, limit int, search string, role Role) ([]User, int64, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
}

type UserService interface {
	GetByID(id uuid.UUID) (*User, error)
	RegisterEmployee(user *User) (*User, error)
	GetByPhoneNumber(phoneNumber string) (*User, error)
	ListUsers(offset, limit int, search string, role Role) ([]User, int64, error)
	UpdateUser(user *User) error
}
