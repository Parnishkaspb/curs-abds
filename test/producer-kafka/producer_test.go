package test

import (
	"github.com/Parnishkaspb/curs-abds/internal/kafka"
	"testing"
)

type KafkaProducer interface {
	SendMessages(repeats int) error
}

type Service struct {
	producer KafkaProducer
}

func New(producer KafkaProducer) *Service {
	return &Service{producer: producer}
}

func TestService(t *testing.T) {
	producer := kafka.NewProducer("localhost", "9092", "transactions")
	defer producer.Close()
	service := New(producer)
	if service.producer == nil {
		t.Fatalf("Producer is nil - это ошибка!")
	}

	err := service.producer.SendMessages(50)
	if err != nil {
		t.Fatalf("Ошибка: %s", err)
	}
}
