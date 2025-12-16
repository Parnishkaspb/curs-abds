.PHONY: mockdata

mockdata:
	@echo "Генерация тестовых данных..."
	DB_DSN=$(DB_DSN) go run ./cmd/mockdata
	@echo "Генерация тестовых закончена"

.PHONY: help

help:
	@echo "Скопируйте команды ниже, чтобы настроить goose"
	@echo "  export GOOSE_DRIVER=postgres"
	@echo "  export GOOSE_DBSTRING=\"host=localhost port=5432 user=user password=password dbname=abds sslmode=disable\""
	@echo "  export GOOSE_MIGRATION_DIR=./internal/database/migrations"
