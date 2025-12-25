package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=user password=password dbname=abds port=5432 sslmode=disable"
	var err error

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("невозможно подключение к базе данных. Ошибка: %s", err)
	}

	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Ошибка при создании миграции. Ошибка: %s", err)
	}
}
