package repository

import (
	"context"
	"time"

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
	return a.client.Get(ctx, key).Result()
}

func (a *Auth) Delete(ctx context.Context, key string) error {
	return a.client.Del(ctx, key).Err()
}
