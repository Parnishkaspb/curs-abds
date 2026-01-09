package service

import (
	"encoding/json"
	"github.com/Parnishkaspb/curs-abds/internal/service/models"
	"github.com/jackc/pgx/v5/pgtype"
)

import (
	"context"
	"errors"

	"github.com/Parnishkaspb/curs-abds/internal/kafka"
)

func (s *DBService) CreateTransactionFromRequest(ctx context.Context, req kafka.TransactionRequest, source, status uint64) (uint64, error) {
	if req.Country == "" {
		return 0, errors.New("country is empty")
	}

	countries, err := s.GetCountries(ctx)
	if err != nil {
		return 0, err
	}

	idCountry, exists := countries[req.Country]
	if !exists {
		idCountry, err = s.SetCountry(ctx, req.Country)
		if err != nil {
			return 0, err
		}
	}

	currencies, err := s.GetCurrencies(ctx)
	if err != nil {
		return 0, err
	}

	idCurrency, exists := currencies[req.Country]
	if !exists {
		return 0, errors.New("currency not found")
	}

	payloadJSON, err := json.Marshal(req)
	if err != nil {
		return 0, err
	}

	newTransaction := models.Transaction{
		TransactionID: pgtype.UUID{},
		AccountID:     req.AccountID,
		Amount:        req.Amount,
		CurrencyID:    idCurrency,
		Merchant:      req.Merchant,
		CountryID:     idCountry,
		StatusID:      status,
		Payload:       string(payloadJSON),
		SourceID:      source,
		CreatedAt:     req.CreatedAt,
	}

	return s.SetTransaction(ctx, newTransaction)
}
