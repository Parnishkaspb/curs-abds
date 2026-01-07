package redisrepo

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type CountryRepository interface {
	SetCountry(ctx context.Context, key, country string, ttl time.Duration) error
	GetCountry(ctx context.Context, key string) (string, error)
}

type RedisCountryRepository struct {
	client *redis.Client
}

func NewRedisCountryRepository(client *redis.Client) *RedisCountryRepository {
	return &RedisCountryRepository{client: client}
}

func (r *RedisCountryRepository) SetCountry(ctx context.Context, key, country string, ttl time.Duration) error {
	return r.client.Set(ctx, key, country, ttl).Err()
}

func (r *RedisCountryRepository) GetCountry(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}
