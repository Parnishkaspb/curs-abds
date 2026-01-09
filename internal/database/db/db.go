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
