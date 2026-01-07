package clickhouse

import (
	"github.com/Parnishkaspb/curs-abds/internal/kafka"
)

type ClickHouseRepo interface {
	AddAcceptedTransaction(req kafka.TransactionRequest) error
	AddDeclineTransaction(req kafka.TransactionRequest) error
}

type ClickHouseService struct {
	repo ClickHouseRepo
}

func NewClickService(repo ClickHouseRepo) *ClickHouseService {
	return &ClickHouseService{
		repo: repo,
	}
}

func (s *ClickHouseService) AddAcceptedTransaction(req kafka.TransactionRequest) error {
	if err := s.repo.AddAcceptedTransaction(req); err != nil {
		return err
	}

	return nil
}

func (s *ClickHouseService) AddDeclineTransaction(req kafka.TransactionRequest) error {
	if err := s.repo.AddDeclineTransaction(req); err != nil {
		return err
	}

	return nil
}
