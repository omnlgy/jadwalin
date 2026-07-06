package repository

import (
	"context"
	"errors"
	"time"

	"github.com/omnlgy/jadwalin/internal/domain"
	"github.com/redis/go-redis/v9"
)

type Auth struct {
	client *redis.Client
}

func NewAuthRepository(client *redis.Client) *Auth {
	return &Auth{
		client: client,
	}
}

func (a *Auth) Create(ctx context.Context, key string, value string, duration time.Duration) error {
	a.client.Set(ctx, key, value, duration)
	return nil
}

func (a *Auth) Get(ctx context.Context, key string) (string, error) {
	result, err := a.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", domain.ErrNotFound
	}
	return result, err
}

func (a *Auth) Delete(ctx context.Context, key string) error {
	return a.client.Del(ctx, key).Err()
}
