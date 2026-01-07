package frauds

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Parnishkaspb/curs-abds/internal/kafka"
	"github.com/Parnishkaspb/curs-abds/internal/service"
	service2 "github.com/Parnishkaspb/curs-abds/internal/service/redis"
)

type Frauds struct {
	rules          []service.EnableFraudRule
	countryService *service2.CountryService
}

func NewFrauds(countryService *service2.CountryService) *Frauds {
	// TODO: в будущем подгрузить из БД
	defaultRules := []service.EnableFraudRule{
		{Code: "high_amount", Threshold: 500000, Severity: "HIGH"}, // 500.000 у.е.
		{Code: "geo_jump", Threshold: 60, Severity: "LOW"},         // 60 секунд TTL
	}

	return &Frauds{
		rules:          defaultRules,
		countryService: countryService,
	}
}

func (f *Frauds) CheckMessage(message string) {
	ctx := context.Background()

	var request kafka.TransactionRequest
	if err := json.Unmarshal([]byte(message), &request); err != nil {
		log.Printf("Error unmarshalling frauds request: %v", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(f.rules))

	for _, rule := range f.rules {
		rule := rule

		go func() {
			defer wg.Done()

			switch rule.Code {
			case "high_amount":
				f.checkHighAmount(rule, request)

			case "geo_jump":
				f.checkGeoJump(ctx, rule, request)
			}
		}()
	}

	wg.Wait()
}

func (f *Frauds) Process(msg []byte) {
	f.CheckMessage(string(msg))
}

// --------------------- FRAUD RULES ----------------------------

func (f *Frauds) checkHighAmount(rule service.EnableFraudRule, req kafka.TransactionRequest) {
	if req.Amount >= rule.Threshold {
		log.Printf(
			"[HIGH_AMOUNT] Сумма(%d) >= лимита(%d). Severity=%s",
			req.Amount, rule.Threshold, rule.Severity,
		)
	}
}

func (f *Frauds) checkGeoJump(ctx context.Context, rule service.EnableFraudRule, req kafka.TransactionRequest) {

	lastCountry, err := f.countryService.GetLastCountry(ctx, req.AccountID)
	if err != nil && !errors.Is(err, service2.ErrCountryNotFound) {
		log.Printf("Error getting last country: %v", err)
		return
	}

	ttl := time.Duration(rule.Threshold) * time.Second
	if err := f.countryService.SaveLastCountry(ctx, req.AccountID, req.Country, ttl); err != nil {
		log.Printf("Error saving country: %v", err)
		return
	}

	if errors.Is(err, service2.ErrCountryNotFound) {
		return
	}

	if lastCountry != req.Country {
		log.Printf(
			"[GEO_JUMP] Аккаунт %d прыгнул из %s в %s. Severity=%s",
			req.AccountID, lastCountry, req.Country, rule.Severity,
		)
	}
}
