package service

import (
	"context"
	"errors"
	"fmt"
	redisrepo "github.com/Parnishkaspb/curs-abds/internal/database/redis"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrCountryNotFound = errors.New("country not found")

type CountryService struct {
	repo redisrepo.CountryRepository
}

func NewCountryService(repo redisrepo.CountryRepository) *CountryService {
	return &CountryService{repo: repo}
}

func (s *CountryService) SaveLastCountry(ctx context.Context, accountID uint64, county string, ttl time.Duration) error {
	if county == "" {
		return errors.New("invalid county")
	}

	key := fmt.Sprintf("fraud:country:%d", accountID)

	return s.repo.SetCountry(ctx, key, county, ttl)
}

func (s *CountryService) GetLastCountry(ctx context.Context, accountID uint64) (string, error) {
	key := fmt.Sprintf("fraud:country:%d", accountID)

	val, err := s.repo.GetCountry(ctx, key)

	if errors.Is(err, redis.Nil) {
		return "", ErrCountryNotFound
	}

	if err != nil {
		return "", err
	}

	return val, nil
}
