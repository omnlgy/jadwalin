package service

import (
	"github.com/google/uuid"
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

func (s *UserService) GetByID(id uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) GetByPhoneNumber(phone string) (*domain.User, error) {
	return s.userRepo.GetByPhoneNumber(phone)
}

func (s *UserService) ListUsers(offset, limit int, search string, role domain.Role) ([]domain.User, int64, error) {
	return s.userRepo.List(offset, limit, search, role)
}

func (s *UserService) UpdateUser(user *domain.User) error {
	return s.userRepo.Update(user)
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}
