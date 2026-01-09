package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Parnishkaspb/curs-abds/internal/service/models"
)

type DatabaseRepo interface {
	GetCountries(ctx context.Context) ([]models.Country, error)
	SetCountry(ctx context.Context, name string) (uint64, error)
	GetCurrencies(ctx context.Context) ([]models.Currency, error)
	SetTransaction(ctx context.Context, req models.Transaction) (uint64, error)
}

type DBService struct {
	repo DatabaseRepo
}

func New(repo DatabaseRepo) *DBService {
	return &DBService{
		repo: repo,
	}
}

func (s *DBService) GetCountries(ctx context.Context) (map[string]uint64, error) {
	result := make(map[string]uint64)

	countries, err := s.repo.GetCountries(ctx)
	if err != nil {
		return nil, err
	}

	for _, country := range countries {
		result[country.Name] = country.ID
	}

	return result, nil
}

func (s *DBService) SetCountry(ctx context.Context, name string) (uint64, error) {
	if name == "" {
		return 0, errors.New("empty country name")
	}

	return s.repo.SetCountry(ctx, strings.ToUpper(name))
}

func (s *DBService) GetCurrencies(ctx context.Context) (map[string]uint64, error) {
	result := make(map[string]uint64)

	currencies, err := s.repo.GetCurrencies(ctx)
	if err != nil {
		return nil, err
	}

	for _, currency := range currencies {
		result[currency.ISO] = currency.ID
	}

	return result, nil
}

func (s *DBService) SetTransaction(ctx context.Context, req models.Transaction) (uint64, error) {
	return s.repo.SetTransaction(ctx, req)
}
