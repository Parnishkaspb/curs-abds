package main

import (
	"context"
	"github.com/Parnishkaspb/curs-abds/internal/kafka"
)

func main() {
	producer := kafka.NewProducer("localhost", "29092", "transactions")
	defer producer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go producer.Start(ctx, 100)

	select {}
}
