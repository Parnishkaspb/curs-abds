package models

import (
	"github.com/jackc/pgx/v5/pgtype"

	"time"
)

type Transaction struct {
	ID            uint64      `gorm:"primaryKey" json:"id"`
	TransactionID pgtype.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"transaction_id"`
	AccountID     uint64      `gorm:"not null;check:account_id >= 0" json:"account_id"`
	Amount        uint64      `gorm:"not null;check:amount >= 0" json:"amount"`
	CurrencyID    uint64      `gorm:"not null" json:"currency_id"`
	Currency      Currency    `gorm:"foreignKey:CurrencyID;references:ID" json:"currency"`
	Merchant      string      `gorm:"type:varchar(255)" json:"merchant"`
	CountryID     uint64      `gorm:"not null" json:"country_id"`
	Country       Country     `gorm:"foreignKey:CountryID;references:ID" json:"country"`
	StatusID      uint64      `gorm:"not null" json:"status_id"`
	Status        Status      `gorm:"foreignKey:StatusID;references:ID" json:"status"`
	Payload       string      `gorm:"type:jsonb" json:"payload"`
	SourceID      uint64      `gorm:"not null" json:"source_id"`
	Source        Source      `gorm:"foreignKey:SourceID;references:ID" json:"source"`
	CreatedAt     time.Time   `gorm:"default:now()" json:"created_at"`
	IngestedAt    time.Time   `json:"ingested_at"`
}

type Currency struct {
	ID   uint64 `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(100)" json:"name"`
	ISO  string `gorm:"type:varchar(50)" json:"iso"`
}

type Country struct {
	ID   uint64 `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(30)" json:"name"`
}

type Status struct {
	ID   uint64 `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(30)" json:"name"`
}

type Source struct {
	ID   uint64 `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(30)" json:"name"`
}

type FraudRule struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"type:varchar(100)" json:"code"`
	Title       string    `gorm:"type:varchar(255)" json:"title"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	Threshold   uint64    `gorm:"not null;check:threshold >= 0" json:"threshold"`
	Enable      bool      `gorm:"default:true" json:"enable"`
	Severity    string    `gorm:"type:varchar(10)" json:"severity"`
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
}

type EnableFraudRule[T uint64 | []string] struct {
	Code      string
	Threshold T
	Severity  string
}

type TransactionFilter struct {
	ID            *uint64
	TransactionID string
	AccountID     *uint64

	StatusID  *uint64
	StatusIDs []uint64
	SourceID  *uint64
	CountryID *uint64
	Merchant  string
	Accepted  *bool

	CreatedFrom *time.Time
	CreatedTo   *time.Time

	Limit   int
	Offset  int
	OrderBy string
}

type TransactionList struct {
	Total  int64
	Limit  int
	Offset int
	Items  []Transaction
}

type TransactionsQuery struct {
	// Идентификаторы
	ID            *uint64 `query:"id"`
	TransactionID string  `query:"transaction_id"`
	AccountID     *uint64 `query:"account_id"`

	// Фильтры
	StatusID  *uint64 `query:"status_id"`
	StatusIDs string  `query:"status_ids"`
	SourceID  *uint64 `query:"source_id"`
	CountryID *uint64 `query:"country_id"`
	Merchant  string  `query:"merchant"`
	Accepted  *bool   `query:"accepted"`

	CreatedFrom string `query:"created_from"`
	CreatedTo   string `query:"created_to"`

	Limit  int    `query:"limit"`
	Offset int    `query:"offset"`
	Sort   string `query:"sort"`
}
