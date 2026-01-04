package consumer_kafka

import (
	kafkaConsumerImport "github.com/Parnishkaspb/curs-abds/internal/kafka"
	"github.com/segmentio/kafka-go"
	"log"
	"testing"
	"time"
)

type KafkaConsumer interface {
	ReadMessages(timeout time.Duration) ([]kafka.Message, error)
}

type Service struct {
	consumer KafkaConsumer
}

func NewService(consumer KafkaConsumer) *Service {
	return &Service{consumer: consumer}
}

func TestService(t *testing.T) {
	kafkaConsumer := kafkaConsumerImport.NewConsumer("localhost", "9092", "transactions")
	service := NewService(kafkaConsumer)

	if service.consumer == nil {
		log.Fatalf("ошибка при создании consumer")
	}

	messages, _ := service.consumer.ReadMessages(2 * time.Second)
	log.Println(messages)

	for _, message := range messages {
		log.Println(message.Key, message.Value)
	}
}
