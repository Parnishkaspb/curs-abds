//package frauds
//
//import (
//	"context"
//	"encoding/json"
//	"errors"
//	"fmt"
//	"github.com/Parnishkaspb/curs-abds/internal/service"
//	"github.com/Parnishkaspb/curs-abds/internal/service/models"
//	"log"
//	"slices"
//	"sync"
//	"time"
//
//	"github.com/Parnishkaspb/curs-abds/internal/kafka"
//	"github.com/Parnishkaspb/curs-abds/internal/service/clickhouse"
//	service2 "github.com/Parnishkaspb/curs-abds/internal/service/redis"
//)
//
//type FraudResult struct {
//	Decline bool
//	Reason  string
//	Log     string
//}
//
//type RuleFunc func(context.Context, kafka.TransactionRequest) FraudResult
//
//type Rule struct {
//	Code string
//	Func RuleFunc
//}
//
//type Frauds struct {
//	rules          []Rule
//	countryService *service2.CountryService
//	clickhouse     *clickhouse.ClickHouseService
//	DBService      *service.DBService
//}
//
//func NewFrauds(
//	countryService *service2.CountryService,
//	clickhouse *clickhouse.ClickHouseService,
//	DBService *service.DBService,
//) *Frauds {
//
//	f := &Frauds{
//		countryService: countryService,
//		clickhouse:     clickhouse,
//		DBService:      DBService,
//	}
//
//	f.rules = []Rule{
//		{
//			Code: "high_amount",
//			Func: f.checkHighAmount,
//		},
//		{
//			Code: "geo_jump",
//			Func: f.checkGeoJump,
//		},
//		{
//			Code: "blacklist",
//			Func: f.blacklist,
//		},
//	}
//
//	return f
//}
//
//func (f *Frauds) Process(msg []byte) {
//	f.CheckMessage(string(msg))
//}
//
//func (f *Frauds) CheckMessage(message string) {
//	ctx := context.Background()
//
//	var req kafka.TransactionRequest
//	if err := json.Unmarshal([]byte(message), &req); err != nil {
//		log.Printf("Error unmarshalling frauds request: %v", err)
//		return
//	}
//
//	results := make(chan FraudResult, len(f.rules))
//	wg := sync.WaitGroup{}
//	wg.Add(len(f.rules))
//
//	for _, rule := range f.rules {
//		rule := rule
//		go func() {
//			defer wg.Done()
//			results <- rule.Func(ctx, req)
//		}()
//	}
//
//	wg.Wait()
//	close(results)
//
//	decline := false
//	declineReason := ""
//	logs := []string{}
//
//	for res := range results {
//		if res.Log != "" {
//			log.Println(res.Log)
//			logs = append(logs, res.Log)
//		}
//
//		if res.Decline {
//			decline = true
//			declineReason = res.Reason
//			break
//		}
//	}
//
//	// Финальное решение
//	if decline {
//		log.Printf("[DECLINE] Причина=%s", declineReason)
//		_ = f.clickhouse.AddDeclineTransaction(req)
//		_, _ = f.DBService.CreateTransactionFromRequest(ctx, req, 1, 2)
//	} else {
//		log.Printf("[ACCEPT]")
//		_ = f.clickhouse.AddAcceptedTransaction(req)
//		_, _ = f.DBService.CreateTransactionFromRequest(ctx, req, 1, 1)
//	}
//
//}
//
////
//// ---------- RULES ----------
////
//
//// RULE 1: HIGH_AMOUNT
//func (f *Frauds) checkHighAmount(ctx context.Context, req kafka.TransactionRequest) FraudResult {
//	rule := models.EnableFraudRule[uint64]{
//		Code:      "high_amount",
//		Threshold: 500000,
//		Severity:  "HIGH",
//	}
//
//	if req.Amount >= rule.Threshold {
//		msg := "[HIGH_AMOUNT] Сумма(%d) >= лимита(%d). Severity=%s"
//		logLine := sprintf(msg, req.Amount, rule.Threshold, rule.Severity)
//		return FraudResult{
//			Decline: true,
//			Reason:  "high_amount",
//			Log:     logLine,
//		}
//	}
//
//	return FraudResult{
//		Decline: false,
//	}
//}
//
//// RULE 2: GEO_JUMP
//func (f *Frauds) checkGeoJump(ctx context.Context, req kafka.TransactionRequest) FraudResult {
//	rule := models.EnableFraudRule[uint64]{
//		Code:      "geo_jump",
//		Threshold: 60,
//		Severity:  "LOW",
//	}
//
//	lastCountry, err := f.countryService.GetLastCountry(ctx, req.AccountID)
//	if err != nil && !errors.Is(err, service2.ErrCountryNotFound) {
//		log.Printf("Error getting last country: %v", err)
//		return FraudResult{}
//	}
//
//	ttl := time.Duration(rule.Threshold) * time.Second
//	if err := f.countryService.SaveLastCountry(ctx, req.AccountID, req.Country, ttl); err != nil {
//		log.Printf("Error saving country: %v", err)
//		return FraudResult{}
//	}
//
//	if errors.Is(err, service2.ErrCountryNotFound) {
//		return FraudResult{Decline: false}
//	}
//
//	if lastCountry != req.Country {
//		msg := "[GEO_JUMP] Аккаунт %d прыгнул из %s в %s. Severity=%s"
//		logLine := sprintf(msg, req.AccountID, lastCountry, req.Country, rule.Severity)
//		return FraudResult{
//			Decline: true,
//			Reason:  "geo_jump",
//			Log:     logLine,
//		}
//	}
//
//	return FraudResult{Decline: false}
//}
//
//// RULE 3: BlackList
//func (f *Frauds) blacklist(ctx context.Context, req kafka.TransactionRequest) FraudResult {
//	rule := models.EnableFraudRule[[]string]{
//		Code:      "blacklist",
//		Threshold: []string{"YMARKET", "CAMOKAT"},
//		Severity:  "HIGH",
//	}
//
//	if ok := slices.Contains(rule.Threshold, req.Merchant); ok {
//		msg := "[BLACK_LIST] Черный список merchant! Merchant: %s. Severity=%s"
//		logLine := sprintf(msg, req.Merchant, rule.Severity)
//		return FraudResult{
//			Decline: true,
//			Reason:  "blacklist",
//			Log:     logLine,
//		}
//	}
//
//	return FraudResult{
//		Decline: false,
//	}
//}
//
////
//// helper
////
//
//func sprintf(format string, v ...interface{}) string {
//	return fmt.Sprintf(format, v...)
//}

package frauds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/Parnishkaspb/curs-abds/internal/kafka"
	"github.com/Parnishkaspb/curs-abds/internal/service"
	"github.com/Parnishkaspb/curs-abds/internal/service/clickhouse"
	"github.com/Parnishkaspb/curs-abds/internal/service/models"
	service2 "github.com/Parnishkaspb/curs-abds/internal/service/redis"
)

type FraudResult struct {
	Decline bool
	Reason  string
	Log     string
}

type Decision struct {
	Decline bool
	Reason  string
	Logs    []string
}

type RuleFunc func(context.Context, kafka.TransactionRequest) FraudResult

type Rule struct {
	Code string
	Func RuleFunc
}

type Frauds struct {
	rules              []Rule
	countryService     *service2.CountryService
	clickhouse         *clickhouse.ClickHouseService
	DBService          *service.DBService
	stopOnFirstDecline bool
}

func NewFrauds(
	countryService *service2.CountryService,
	clickhouse *clickhouse.ClickHouseService,
	DBService *service.DBService,
) *Frauds {

	f := &Frauds{
		countryService:     countryService,
		clickhouse:         clickhouse,
		DBService:          DBService,
		stopOnFirstDecline: true,
	}

	f.rules = []Rule{
		{Code: "high_amount", Func: f.checkHighAmount},
		{Code: "geo_jump", Func: f.checkGeoJump},
		{Code: "blacklist", Func: f.blacklist},
	}

	return f
}

// Опционально: хочешь собирать логи/результаты ВСЕХ правил даже если уже decline
func (f *Frauds) SetStopOnFirstDecline(v bool) {
	f.stopOnFirstDecline = v
}

// Evaluate — чистая логика: считает решение по правилам и возвращает Decision.
// Никаких записей в ClickHouse/Postgres тут нет.
func (f *Frauds) Evaluate(ctx context.Context, req kafka.TransactionRequest) Decision {
	results := make(chan FraudResult, len(f.rules))
	wg := sync.WaitGroup{}
	wg.Add(len(f.rules))

	for _, rule := range f.rules {
		rule := rule
		go func() {
			defer wg.Done()
			results <- rule.Func(ctx, req)
		}()
	}

	wg.Wait()
	close(results)

	decision := Decision{
		Decline: false,
		Reason:  "",
		Logs:    make([]string, 0, len(f.rules)),
	}

	for res := range results {
		if res.Log != "" {
			decision.Logs = append(decision.Logs, res.Log)
		}
		if res.Decline && !decision.Decline {
			decision.Decline = true
			decision.Reason = res.Reason
			if f.stopOnFirstDecline {
				break
			}
		}
	}

	return decision
}

func (f *Frauds) ApplyDecision(ctx context.Context, req kafka.TransactionRequest, decision Decision, sourceID uint64) (uint64, error) {
	for _, line := range decision.Logs {
		log.Println(line)
	}

	if decision.Decline {
		log.Printf("[DECLINE] Причина=%s", decision.Reason)
		_ = f.clickhouse.AddDeclineTransaction(req)
		return f.DBService.CreateTransactionFromRequest(ctx, req, sourceID, 2) // statusID=2
	}

	log.Printf("[ACCEPT]")
	_ = f.clickhouse.AddAcceptedTransaction(req)
	return f.DBService.CreateTransactionFromRequest(ctx, req, sourceID, 1) // statusID=1
}

// Process — как раньше: принимает msg bytes (Kafka), внутри парсит и применяет.
func (f *Frauds) Process(msg []byte) {
	ctx := context.Background()

	var req kafka.TransactionRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		log.Printf("Error unmarshalling frauds request: %v", err)
		return
	}

	decision := f.Evaluate(ctx, req)

	_, _ = f.ApplyDecision(ctx, req, decision, 1)
}

func (f *Frauds) CheckMessage(message string) {
	f.Process([]byte(message))
}

//
// ---------- RULES ----------
//

// RULE 1: HIGH_AMOUNT
func (f *Frauds) checkHighAmount(ctx context.Context, req kafka.TransactionRequest) FraudResult {
	rule := models.EnableFraudRule[uint64]{
		Code:      "high_amount",
		Threshold: 500000,
		Severity:  "HIGH",
	}

	if req.Amount >= rule.Threshold {
		msg := "[HIGH_AMOUNT] Сумма(%d) >= лимита(%d). Severity=%s"
		logLine := sprintf(msg, req.Amount, rule.Threshold, rule.Severity)
		return FraudResult{
			Decline: true,
			Reason:  "high_amount",
			Log:     logLine,
		}
	}

	return FraudResult{Decline: false}
}

// RULE 2: GEO_JUMP
func (f *Frauds) checkGeoJump(ctx context.Context, req kafka.TransactionRequest) FraudResult {
	rule := models.EnableFraudRule[uint64]{
		Code:      "geo_jump",
		Threshold: 60,
		Severity:  "LOW",
	}

	lastCountry, err := f.countryService.GetLastCountry(ctx, req.AccountID)
	if err != nil && !errors.Is(err, service2.ErrCountryNotFound) {
		log.Printf("Error getting last country: %v", err)
		return FraudResult{} // нейтрально
	}

	ttl := time.Duration(rule.Threshold) * time.Second
	if err := f.countryService.SaveLastCountry(ctx, req.AccountID, req.Country, ttl); err != nil {
		log.Printf("Error saving country: %v", err)
		return FraudResult{}
	}

	if errors.Is(err, service2.ErrCountryNotFound) {
		return FraudResult{Decline: false}
	}

	if lastCountry != req.Country {
		msg := "[GEO_JUMP] Аккаунт %d прыгнул из %s в %s. Severity=%s"
		logLine := sprintf(msg, req.AccountID, lastCountry, req.Country, rule.Severity)
		return FraudResult{
			Decline: true,
			Reason:  "geo_jump",
			Log:     logLine,
		}
	}

	return FraudResult{Decline: false}
}

// RULE 3: BlackList
func (f *Frauds) blacklist(ctx context.Context, req kafka.TransactionRequest) FraudResult {
	rule := models.EnableFraudRule[[]string]{
		Code:      "blacklist",
		Threshold: []string{"YMARKET", "CAMOKAT"},
		Severity:  "HIGH",
	}

	if slices.Contains(rule.Threshold, req.Merchant) {
		msg := "[BLACK_LIST] Черный список merchant! Merchant: %s. Severity=%s"
		logLine := sprintf(msg, req.Merchant, rule.Severity)
		return FraudResult{
			Decline: true,
			Reason:  "blacklist",
			Log:     logLine,
		}
	}

	return FraudResult{Decline: false}
}

//
// helper
//

func sprintf(format string, v ...interface{}) string {
	return fmt.Sprintf(format, v...)
}
