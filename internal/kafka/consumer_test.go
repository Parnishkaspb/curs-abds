package kafka_test

import (
	"context"
	"errors"
	kafka2 "github.com/Parnishkaspb/curs-abds/internal/kafka"
	"github.com/segmentio/kafka-go"
	"testing"
	"time"
)

//
// ===== MOCK READER =====
//

type MockReader struct {
	Messages     []kafka.Message
	Err          error
	ReadCalls    int
	Closed       bool
	CurrentIndex int
}

func (m *MockReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	m.ReadCalls++

	if m.Err != nil {
		return kafka.Message{}, m.Err
	}
	if m.CurrentIndex >= len(m.Messages) {
		<-ctx.Done()
		return kafka.Message{}, ctx.Err()
	}

	msg := m.Messages[m.CurrentIndex]
	m.CurrentIndex++
	return msg, nil
}

func (m *MockReader) Close() error {
	m.Closed = true
	return nil
}

type MockProcessor struct {
	Calls int
	Last  []byte
}

func (m *MockProcessor) Process(msg []byte) {
	m.Calls++
	m.Last = msg
}

//
// ===== TESTS =====
//

func TestConsumer_Start_ReadsMessages(t *testing.T) {
	mock := &MockReader{
		Messages: []kafka.Message{
			{Key: []byte("1"), Value: []byte("hello")},
			{Key: []byte("2"), Value: []byte("world")},
		},
	}

	proc := &MockProcessor{}
	consumer := kafka2.NewConsumerWithReader(mock, proc)

	ctx, cancel := context.WithCancel(context.Background())
	go consumer.Start(ctx)

	time.Sleep(50 * time.Millisecond)
	cancel()

	if mock.ReadCalls == 0 {
		t.Fatalf("expected ReadMessage to be called")
	}
	if proc.Calls == 0 {
		t.Fatalf("processor should be called at least once")
	}
}

func TestConsumer_Start_ReadError(t *testing.T) {
	mock := &MockReader{
		Err: errors.New("read failed"),
	}
	proc := &MockProcessor{}
	consumer := kafka2.NewConsumerWithReader(mock, proc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go consumer.Start(ctx)
	time.Sleep(50 * time.Millisecond)

	if mock.ReadCalls == 0 {
		t.Fatalf("expected ReadMessage to be called even on error")
	}
	if proc.Calls != 0 {
		t.Fatalf("processor must NOT be called on read error")
	}
}

func TestConsumer_Start_ContextCancelStops(t *testing.T) {
	mock := &MockReader{
		Messages: []kafka.Message{
			{Value: []byte("x")},
		},
	}
	proc := &MockProcessor{}
	consumer := kafka2.NewConsumerWithReader(mock, proc)

	ctx, cancel := context.WithCancel(context.Background())
	go consumer.Start(ctx)

	time.Sleep(30 * time.Millisecond)
	cancel()

	prev := mock.ReadCalls
	time.Sleep(50 * time.Millisecond)

	if mock.ReadCalls != prev {
		t.Fatalf("consumer kept reading after context cancel")
	}
}

func TestConsumer_Close(t *testing.T) {
	mock := &MockReader{}
	proc := &MockProcessor{}

	consumer := kafka2.NewConsumerWithReader(mock, proc)
	consumer.Close()

	if !mock.Closed {
		t.Fatalf("reader.Close() should be called")
	}
}

func TestConsumer_Table(t *testing.T) {
	tests := []struct {
		name        string
		mock        *MockReader
		wantReads   int
		cancelAfter time.Duration
	}{
		{
			name: "reads two messages",
			mock: &MockReader{
				Messages: []kafka.Message{
					{Value: []byte("a")},
					{Value: []byte("b")},
				},
			},
			wantReads:   2,
			cancelAfter: 40 * time.Millisecond,
		},
		{
			name: "read error",
			mock: &MockReader{
				Err: errors.New("failed"),
			},
			wantReads:   1,
			cancelAfter: 20 * time.Millisecond,
		},
		{
			name: "cancel before messages",
			mock: &MockReader{
				Messages: []kafka.Message{},
			},
			wantReads:   0,
			cancelAfter: 5 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proc := &MockProcessor{}
			c := kafka2.NewConsumerWithReader(tt.mock, proc)

			ctx, cancel := context.WithCancel(context.Background())
			go c.Start(ctx)

			time.Sleep(tt.cancelAfter)
			cancel()
			time.Sleep(30 * time.Millisecond)

			if tt.mock.ReadCalls < tt.wantReads {
				t.Fatalf("expected at least %d reads, got %d",
					tt.wantReads, tt.mock.ReadCalls)
			}
		})
	}
}
