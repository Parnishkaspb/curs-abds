package db

import (
	"context"

	"github.com/Parnishkaspb/curs-abds/internal/service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Repo struct {
	db *gorm.DB
}

func New() *Repo {
	dsn := "host=localhost user=user password=password dbname=abds port=5432 sslmode=disable"
	var err error

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("невозможно подключение к базе данных. Ошибка: %s", err)
	}

	return &Repo{db: db}
}

func (r *Repo) GetCountries(ctx context.Context) ([]models.Country, error) {
	var countries []models.Country
	err := r.db.WithContext(ctx).Find(&countries).Error
	return countries, err
}

func (r *Repo) SetCountry(ctx context.Context, name string) (uint64, error) {
	country := models.Country{
		Name: name,
	}

	err := r.db.WithContext(ctx).
		Create(&country).
		Error
	if err != nil {
		return 0, err
	}

	return country.ID, nil
}

func (r *Repo) GetCurrencies(ctx context.Context) ([]models.Currency, error) {
	var currencies []models.Currency
	err := r.db.WithContext(ctx).Find(&currencies).Error
	return currencies, err
}

func (r *Repo) SetTransaction(ctx context.Context, req models.Transaction) (uint64, error) {
	err := r.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return 0, err
	}

	return req.ID, nil
}

func (r *Repo) SearchTransactions(ctx context.Context, f models.TransactionFilter) (models.TransactionList, error) {
	q := r.db.WithContext(ctx).
		Model(&models.Transaction{}).
		Preload("Currency").
		Preload("Country").
		Preload("Status").
		Preload("Source")

	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.TransactionID != "" {
		q = q.Where("transaction_id = ?", f.TransactionID)
	}
	if f.AccountID != nil {
		q = q.Where("account_id = ?", *f.AccountID)
	}
	if f.StatusID != nil {
		q = q.Where("status_id = ?", *f.StatusID)
	}
	if len(f.StatusIDs) > 0 {
		q = q.Where("status_id IN ?", f.StatusIDs)
	}
	if f.SourceID != nil {
		q = q.Where("source_id = ?", *f.SourceID)
	}
	if f.CountryID != nil {
		q = q.Where("country_id = ?", *f.CountryID)
	}
	if f.Merchant != "" {
		q = q.Where("merchant ILIKE ?", "%"+f.Merchant+"%")
	}
	if f.CreatedFrom != nil {
		q = q.Where("created_at >= ?", *f.CreatedFrom)
	}
	if f.CreatedTo != nil {
		q = q.Where("created_at <= ?", *f.CreatedTo)
	}

	if f.Accepted != nil {
		q = q.Where("accepted = ?", *f.Accepted)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return models.TransactionList{}, err
	}

	var items []models.Transaction
	if err := q.Order(f.OrderBy).Limit(f.Limit).Offset(f.Offset).Find(&items).Error; err != nil {
		return models.TransactionList{}, err
	}

	return models.TransactionList{
		Total:  total,
		Limit:  f.Limit,
		Offset: f.Offset,
		Items:  items,
	}, nil
}
