package service

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/omnlgy/jadwalin/internal/domain"
)

type AuthService struct {
	authRepo domain.AuthRepository
}

func NewAuthService(authRepo domain.AuthRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

func (s *AuthService) GenerateOTP(ctx context.Context, key string) (string, error) {
	code := generateCode()
	if err := s.authRepo.Create(ctx, key, code, 5*time.Minute); err != nil {
		return "", err
	}

	return code, nil
}

func (s *AuthService) VerifyOTP(ctx context.Context, key, code string) error {
	exists, err := s.authRepo.Get(ctx, key)

	if err != nil {
		return err
	}

	if exists != code {
		return domain.ErrInvalidOTP
	}
	return nil
}

func (s *AuthService) Login(ctx context.Context, phone string) (string, error) {
	// TODO: implement login logic
	return "", nil
}

func generateCode() string {
	return fmt.Sprintf("%06d", rand.IntN(1000000))
}
