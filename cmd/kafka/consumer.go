package main

import (
	"context"
	redisrepo "github.com/Parnishkaspb/curs-abds/internal/database/redis"
	"github.com/Parnishkaspb/curs-abds/internal/kafka"
	"github.com/Parnishkaspb/curs-abds/internal/service/frauds"
	service "github.com/Parnishkaspb/curs-abds/internal/service/redis"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	repo := redisrepo.NewRedisCountryRepository(rdb)

	countryService := service.NewCountryService(repo)

	fraudProcessor := frauds.NewFrauds(countryService)

	consumer := kafka.NewConsumer(
		"localhost",
		"29092",
		"transactions",
		fraudProcessor,
	)

	consumer.Start(ctx)
}
