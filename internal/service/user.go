package service

import (
	"github.com/gofrs/uuid/v5"
	"github.com/omnlgy/jadwalin/internal/domain"
)

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) RegisterEmployee(user *domain.User) (*domain.User, error) {
	user.Role = domain.RoleEmployee
	return s.userRepo.Create(user)
}

func (s *UserService) VerifyUser(id uuid.UUID) (*domain.User, error) {
	// return s.userRepo.Verify(id)
}

func (s *UserService) GetByID(id uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}
