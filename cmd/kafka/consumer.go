package main

import (
	"context"
	clickhouserepo "github.com/Parnishkaspb/curs-abds/internal/database/clickhouse"
	redisrepo "github.com/Parnishkaspb/curs-abds/internal/database/redis"
	"github.com/Parnishkaspb/curs-abds/internal/kafka"
	"github.com/Parnishkaspb/curs-abds/internal/service/clickhouse"
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

	repoClick := clickhouserepo.NewClickHouse(context.Background())
	clickService := clickhouse.NewClickService(repoClick)

	fraudProcessor := frauds.NewFrauds(countryService, clickService)

	consumer := kafka.NewConsumer(
		"localhost",
		"29092",
		"transactions",
		fraudProcessor,
	)

	consumer.Start(ctx)
}
