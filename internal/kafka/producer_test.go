package kafka_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"testing"
	"time"

	kafka2 "github.com/Parnishkaspb/curs-abds/internal/kafka"
)

//
// ===== MOCK WRITER =====
//

type MockWriter struct {
	MessagesSent []kafka.Message
	Err          error
	Closed       bool
}

func (m *MockWriter) WriteMessages(_ context.Context, msgs ...kafka.Message) error {
	if m.Err != nil {
		return m.Err
	}
	m.MessagesSent = append(m.MessagesSent, msgs...)
	return nil
}

func (m *MockWriter) Close() error {
	m.Closed = true
	return nil
}

//
// ===== OVERRIDDEN CONSTRUCTOR =====
//

//func NewProducerWithWriter(w Writer) *Producer {
//	return &Producer{writer: w}
//}

//
// ===== TESTS =====
//

func TestProducer_SendMessages_Success(t *testing.T) {
	mockWriter := &MockWriter{}
	p := kafka2.NewProducerWithWriter(mockWriter)

	err := p.SendMessages(3)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if len(mockWriter.MessagesSent) != 3 {
		t.Fatalf("expected 3 messages sent, got %d", len(mockWriter.MessagesSent))
	}

	for _, msg := range mockWriter.MessagesSent {
		if len(msg.Value) == 0 {
			t.Fatal("empty message value")
		}
		var body kafka2.TransactionRequest
		if json.Unmarshal(msg.Value, &body) != nil {
			t.Fatal("message value is not valid JSON")
		}
	}
}

func TestProducer_SendMessages_WriteError(t *testing.T) {
	mockWriter := &MockWriter{
		Err: errors.New("write failed"),
	}
	p := kafka2.NewProducerWithWriter(mockWriter)

	err := p.SendMessages(5)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestProducer_CreateMessages(t *testing.T) {
	p := kafka2.NewProducerWithWriter(&MockWriter{})

	msgs, err := p.CreateMessages(4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(msgs) != 4 {
		t.Fatalf("expected 4 messages, got %d", len(msgs))
	}
}

func TestProducer_CreateMessage(t *testing.T) {
	p := kafka2.NewProducerWithWriter(&MockWriter{})

	msg := p.CreateMessage()

	if msg.TransactionID == "" {
		t.Fatal("TransactionID should not be empty")
	}
	if msg.Country == "" {
		t.Fatal("Country should not be empty")
	}
	if msg.Merchant == "" {
		t.Fatal("Merchant should not be empty")
	}
}

func TestProducer_GenerateTransactionID(t *testing.T) {
	p := kafka2.NewProducerWithWriter(&MockWriter{})

	id := p.GenerateTransactionID()

	if len(id) != 10 {
		t.Fatalf("expected length 10, got %d", len(id))
	}
}

func TestProducer_GenerateCountry(t *testing.T) {
	p := kafka2.NewProducerWithWriter(&MockWriter{})

	c := p.GenerateCountryForTransactionRequest()
	if c == "" {
		t.Fatal("country cannot be empty")
	}
}

func TestProducer_GenerateMerchant(t *testing.T) {
	p := kafka2.NewProducerWithWriter(&MockWriter{})

	m := p.GenerateMerchantForTransactionRequest()
	if m == "" {
		t.Fatal("merchant cannot be empty")
	}
}

func TestProducer_Close(t *testing.T) {
	mockWriter := &MockWriter{}
	p := kafka2.NewProducerWithWriter(mockWriter)

	p.Close()

	if !mockWriter.Closed {
		t.Fatal("writer.Close() must be called")
	}
}

//
// ===== TEST Start() =====
//

func TestProducer_Start(t *testing.T) {
	mockWriter := &MockWriter{}
	p := kafka2.NewProducerWithWriter(mockWriter)

	ctx, cancel := context.WithCancel(context.Background())

	go p.Start(ctx, 2)

	time.Sleep(120 * time.Millisecond)
	cancel()

	if len(mockWriter.MessagesSent) == 0 {
		t.Fatal("expected Start to send messages at least once")
	}
}

//
// ===== TABLE TESTS =====
//

func TestProducer_Table(t *testing.T) {
	tests := []struct {
		name    string
		mock    *MockWriter
		repeats int
		wantErr bool
	}{
		{
			name:    "ok single batch",
			mock:    &MockWriter{},
			repeats: 2,
			wantErr: false,
		},
		{
			name:    "write error",
			mock:    &MockWriter{Err: errors.New("cannot write")},
			repeats: 3,
			wantErr: true,
		},
		{
			name:    "zero repeats",
			mock:    &MockWriter{},
			repeats: 0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := kafka2.NewProducerWithWriter(tt.mock)

			err := p.SendMessages(tt.repeats)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}
