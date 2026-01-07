package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

type Reader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

type MessageProcessor interface {
	Process(msg []byte)
}

type Consumer struct {
	reader    Reader
	processor MessageProcessor
}

func NewConsumer(host, port, topic string, processor MessageProcessor) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:         []string{host + ":" + port},
			Topic:           topic,
			GroupID:         "transactions-group",
			MinBytes:        1,
			MaxBytes:        10e6,
			StartOffset:     kafka.FirstOffset,
			ReadLagInterval: -1,
		}),
		processor: processor,
	}
}

func NewConsumerWithReader(r Reader, processor MessageProcessor) *Consumer {
	return &Consumer{
		reader:    r,
		processor: processor,
	}
}

func (c *Consumer) Close() {
	if c.reader != nil {
		c.reader.Close()
	}
}

func (c *Consumer) Start(ctx context.Context) {
	fmt.Println("consumer started...")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("consumer stopped")
			return

		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					fmt.Println("consumer context canceled")
					return
				}

				fmt.Println("read error:", err)
				time.Sleep(time.Second)
				continue
			}

			if c.processor != nil {
				go c.processor.Process(msg.Value)
			}
		}
	}
}
