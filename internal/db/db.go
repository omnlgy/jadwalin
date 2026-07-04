package db

import (
	"github.com/omnlgy/jadwalin/internal/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Redis *redis.Client

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := "host=" + cfg.DB_HOST + " user=" + cfg.DB_USER + " password=" + cfg.DB_PASS + " dbname=" + cfg.DB_NAME + " port=" + cfg.DB_PORT + " sslmode=disable"

	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return DB, nil
}

func NewRedisClient(cfg *config.Config) *redis.Client {
	Redis = redis.NewClient(&redis.Options{
		Addr:     cfg.REDIS_ADDR,
		Password: cfg.REDIS_PASS,
		DB:       0,
	})

	return Redis
}
