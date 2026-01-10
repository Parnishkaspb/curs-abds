package main

import (
	"context"

	clickhouserepo "github.com/Parnishkaspb/curs-abds/internal/database/clickhouse"
	database "github.com/Parnishkaspb/curs-abds/internal/database/db"
	redisrepo "github.com/Parnishkaspb/curs-abds/internal/database/redis"
	"github.com/Parnishkaspb/curs-abds/internal/service"
	"github.com/Parnishkaspb/curs-abds/internal/service/clickhouse"
	frauds2 "github.com/Parnishkaspb/curs-abds/internal/service/frauds"
	"github.com/Parnishkaspb/curs-abds/internal/service/handle"
	service2 "github.com/Parnishkaspb/curs-abds/internal/service/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	repo := database.New()
	svc := service.New(repo)

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	redisRepo := redisrepo.NewRedisCountryRepository(rdb)
	countryService := service2.NewCountryService(redisRepo)

	repoClick := clickhouserepo.NewClickHouse(context.Background())
	clickService := clickhouse.NewClickService(repoClick)

	frauds := frauds2.NewFrauds(countryService, clickService, svc)

	h := handle.NewHandle(svc, frauds)

	e := echo.New()
	e.Use(middleware.Logger())

	e.POST("/transactions", h.CreateTransaction)
	e.GET("/transactions", h.GetTransactions)

	e.Start("localhost:8080")
}
