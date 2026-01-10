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

.PHONY: start
start:
	go run cmd/curs-abds/main.go

.PHONY: start_docker
start_docker:
	cd docker && docker-compose up -d

.PHONY: topic
topic:
	docker exec -it kafka kafka-topics --create --topic transactions --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1


.PHONY: read_kafka_consumer
# Прочитать сообщения
read_kafka_consumer:
	docker exec -it kafka kafka-console-consumer --topic transactions --from-beginning --bootstrap-server localhost:9092 --max-messages $(or $(MESSAGES),100)

.PHONY: start_producer
start_producer:
	go run cmd/kafka/producer.go

.PHONY: start_consumer
start_consumer:
	go run cmd/kafka/consumer.go

# Проверить список топиков
#docker exec -it kafka kafka-topics --list --bootstrap-server localhost:9092

# Отправить тестовое сообщение
#docker exec -it kafka bash -c "echo 'test message' | kafka-console-producer --topic transactions --bootstrap-server localhost:9092"
