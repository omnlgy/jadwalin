package service

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/omnlgy/jadwalin/internal/domain"
	"github.com/omnlgy/jadwalin/utils"
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

func (s *AuthService) Logout(ctx context.Context, token string) error {
	claim, err := utils.ValidateJWT(token)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("blacklist:%s", token)
	err = s.authRepo.Create(ctx, key, "", time.Until(claim.ExpiresAt.Time))

	if err != nil {
		return err
	}
	return nil
}

func generateCode() string {
	return fmt.Sprintf("%06d", rand.IntN(1000000))
}
