package domain

import (
	"context"
	"time"
)

type AuthRepository interface {
	Create(ctx context.Context, key string, value string, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type AuthService interface {
	GenerateOTP(ctx context.Context, key string) (string, error)
	VerifyOTP(ctx context.Context, key, code string) error
}
