package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"math/rand"
	"strings"
	"time"
)

type TransactionRequest struct {
	TransactionID string    `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
	AccountID     uint64    `json:"account_id"`
	Amount        uint64    `json:"amount"`
	Country       string    `json:"country"`
	Merchant      string    `json:"merchant"`
}

type Writer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

//type Producer struct {
//	writer *kafka.Writer
//}

type Producer struct {
	writer Writer
}

func NewProducerWithWriter(w Writer) *Producer {
	return &Producer{writer: w}
}

// New создает нового продюсера
func NewProducer(host, port, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(host + ":" + port),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			BatchSize:    100,
			RequiredAcks: kafka.RequireOne,
			Compression:  kafka.Snappy,
		},
	}
}

// SendMessages отправляет batch сообщений в Kafka
func (p *Producer) SendMessages(repeats int) error {
	messages, err := p.CreateMessages(repeats)
	if err != nil {
		return err
	}

	kafkaMessages := make([]kafka.Message, len(messages))
	for i, msg := range messages {
		jsonData, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message %d: %w", i, err)
		}

		kafkaMessages[i] = kafka.Message{
			Key:   []byte(fmt.Sprintf("%d", msg.AccountID)),
			Value: jsonData,
			Headers: []kafka.Header{
				{Key: "country", Value: []byte(msg.Country)},
				{Key: "merchant", Value: []byte(msg.Merchant)},
			},
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := p.writer.WriteMessages(ctx, kafkaMessages...); err != nil {
		return fmt.Errorf("failed to write messages: %w", err)
	}

	return nil
}

func (p *Producer) CreateMessages(repeats int) ([]TransactionRequest, error) {
	tR := make([]TransactionRequest, repeats)
	for i := 0; i < repeats; i++ {
		tR[i] = p.CreateMessage()
	}
	return tR, nil
}

func (p *Producer) CreateMessage() TransactionRequest {
	return TransactionRequest{
		TransactionID: p.GenerateTransactionID(),
		CreatedAt:     time.Now(),
		AccountID:     rand.Uint64() % 10,
		Amount:        rand.Uint64() % 1000000,
		Country:       p.GenerateCountryForTransactionRequest(),
		Merchant:      p.GenerateMerchantForTransactionRequest(),
	}
}

func (p *Producer) GenerateTransactionID() string {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	var b strings.Builder
	b.Grow(10)
	for i := 0; i < 10; i++ {
		b.WriteByte(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func (p *Producer) GenerateCountryForTransactionRequest() string {
	countries := []string{"RU", "BY", "USA", "ENG", "DECH"}
	return countries[rand.Intn(len(countries))]
}

func (p *Producer) GenerateMerchantForTransactionRequest() string {
	countries := []string{"MEGAMARKET", "OZON", "WILDBERRIES", "YMARKET", "CAMOKAT", "YLAVKA", "CDEK"}
	return countries[rand.Intn(len(countries))]
}

func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}

func (p *Producer) Start(ctx context.Context, repeats int) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("producer stopped")
			return
		case <-ticker.C:
			if err := p.SendMessages(repeats); err != nil {
				fmt.Printf("failed to send messages: %v\n", err)
			} else {
				fmt.Println("messages sent")
			}
		}
	}
}
