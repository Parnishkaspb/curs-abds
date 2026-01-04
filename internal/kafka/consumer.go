package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(host, port, topic string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{host + ":" + port},
			Topic:   topic,
			//GroupID: "my-groupID",
		}),
	}
}

func (c *Consumer) Close() {
	c.reader.Close()
}

func (c *Consumer) ReadMessages(timeout time.Duration) ([]kafka.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var messages []kafka.Message

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			break
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
