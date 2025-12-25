package service

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
