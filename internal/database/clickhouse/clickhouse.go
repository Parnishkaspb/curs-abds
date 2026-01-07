package clickhouserepo

import (
	"context"
	"fmt"
	"github.com/Parnishkaspb/curs-abds/internal/kafka"
	"log"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouseRepo struct {
	db  clickhouse.Conn
	ctx context.Context
}

func NewClickHouse(context context.Context) *ClickHouseRepo {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "abds",
			Username: "user",
			Password: "",
		},
		Debugf: func(format string, v ...interface{}) {
			fmt.Printf(format, v...)
		},
	})

	if err != nil {
		log.Fatalf("Ошибка подключения к ClickHouse: %v", err)
	}

	return &ClickHouseRepo{
		db:  conn,
		ctx: context,
	}
}

func (c *ClickHouseRepo) AddAcceptedTransaction(req kafka.TransactionRequest) error {
	return c.insertTransaction(req, true)
}

func (c *ClickHouseRepo) AddDeclineTransaction(req kafka.TransactionRequest) error {
	return c.insertTransaction(req, false)
}

func (c *ClickHouseRepo) insertTransaction(req kafka.TransactionRequest, accepted bool) error {
	query := `
		INSERT INTO transactions 
		(transaction_id, created_at, account_id, amount, country, merchant, accepted)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	return c.db.Exec(c.ctx, query,
		req.TransactionID,
		req.CreatedAt,
		req.AccountID,
		req.Amount,
		req.Country,
		req.Merchant,
		accepted,
	)
}

func (c *ClickHouseRepo) AddTransactionsBatch(reqs []kafka.TransactionRequest, accepted bool) error {
	batch, err := c.db.PrepareBatch(c.ctx,
		`INSERT INTO transactions 
        (transaction_id, created_at, account_id, amount, country, merchant, accepted)`,
	)
	if err != nil {
		return err
	}

	for _, req := range reqs {
		err = batch.Append(
			req.TransactionID,
			req.CreatedAt,
			req.AccountID,
			req.Amount,
			req.Country,
			req.Merchant,
			accepted,
		)
		if err != nil {
			return err
		}
	}

	return batch.Send()
}
