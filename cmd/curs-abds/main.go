package main

import (
	"encoding/json"
	"github.com/Parnishkaspb/curs-abds/internal/kafka"
	"github.com/Parnishkaspb/curs-abds/internal/service"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

//type TransactionRequest struct {
//	TransactionID string    `json:"transaction_id"`
//	CreatedAt     time.Time `json:"created_at"`
//	AccountID     uint64    `json:"account_id"`
//	Amount        uint64    `json:"amount"`
//	Country       string    `json:"country"`
//	Merchant      string    `json:"merchant"`
//}

var db *gorm.DB

var Countries = map[string]int{
	"RU":  1,
	"BY":  2,
	"USA": 3,
}

var Currencies = map[string]uint64{
	"RU":  1,
	"BY":  2,
	"USA": 3,
}

var transactions []service.Transaction

func initDB() {
	dsn := "host=localhost user=user password=password dbname=abds port=5432 sslmode=disable"
	var err error

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("невозможно подключение к базе данных. Ошибка: %s", err)
	}

	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Ошибка при создании миграции. Ошибка: %s", err)
	}
}

func createTransaction(c echo.Context) error {
	var req kafka.TransactionRequest

	if err := c.Bind(&req); err != nil {
		log.Printf("Ошибка: %s", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Неверный запрос"})
	}

	idCountry, exists := Countries[req.Country]
	if !exists {
		log.Printf("Страна не найдена: %s. Автоматическое добавление!", req.Country)
		idCountry = maxID() + 1
		Countries[req.Country] = idCountry
	}

	idCurrency, exists := Currencies[req.Country]
	if !exists {
		log.Fatalf("Валюта не найдена: %s. Автоматическое добавление невозможно!", req.Country)
	}

	payloadJSON, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("Ошибка при конвертации данных!")
	}

	newTransaction := service.Transaction{
		ID:            maxIDTransaction() + 1,
		TransactionID: pgtype.UUID{},
		AccountID:     req.AccountID,
		Amount:        req.Amount,
		CurrencyID:    idCurrency,
		Merchant:      req.Merchant,
		CountryID:     uint64(idCountry),
		StatusID:      1,
		Payload:       string(payloadJSON),
		SourceID:      1,
		CreatedAt:     req.CreatedAt,
	}

	transactions = append(transactions, newTransaction)

	return c.JSON(http.StatusOK, map[string]string{"message": "Успешное добавление транзакции"})

}

func main() {
	//e := echo.New()
	//e.Use(middleware.Logger())
	////e.Use(fraudsMiddleware())
	//
	//e.GET("/", func(c echo.Context) error {
	//	return c.String(http.StatusOK, "Hello, World!")
	//})
	//
	////e.POST("/transactions", func(c echo.Context) error {
	////	return c.String(http.StatusOK, "Hello, World!")
	////})
	//
	//e.POST("/transactions", createTransaction)

	//rdb := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "", // no password
	//	DB:       0,  // use default DB
	//	Protocol: 2,
	//})
	//
	//ctx := context.Background()
	//
	//res1, err := rdb.ZAdd(ctx, "racer_scores",
	//	redis.Z{Member: "Norem", Score: 10},
	//).Result()
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(res1)

	//e.Start("localhost:8080")
}

func maxID() int {
	maxId := 0
	for _, id := range Countries {
		if id > maxId {
			maxId = id
		}
	}

	return maxId
}

func maxIDTransaction() uint64 {
	var transactionsMap []service.Transaction
	if len(transactionsMap) == 0 {
		return 0
	}
	return transactionsMap[len(transactionsMap)-1].ID
}

var ip = make(map[string]uint)

var frauds = []service.EnableFraudRule{
	{
		Code:      "high_amount",
		Threshold: 10000,
		Severity:  "HIGH",
	},
	{
		Code:      "geo_jump",
		Threshold: 60,
		Severity:  "LOW",
	},
}

//func fraudsMiddleware() echo.MiddlewareFunc {
//	return func(next echo.HandlerFunc) echo.HandlerFunc {
//		return func(c echo.Context) error {
//
//			wg := sync.WaitGroup{}
//			wg.Add(len(frauds))
//
//			for fraud := range frauds {
//				go func(fraud service.EnableFraudRule) error {
//					defer wg.Done()
//
//					switch fraud.Code {
//					case "high_amount":
//						var req TransactionRequest
//
//						if err := c.Bind(&req); err != nil {
//							log.Printf("Ошибка: %s", err.Error())
//						}
//
//						if req.Amount >= fraud.Threshold {
//							return c.JSON(http.StatusConflict, map[string]string{
//								"message": "Проблема с суммой перевода",
//							})
//						}
//						break
//
//					}
//				}()
//			}
//
//			realIP := c.RealIP()
//			count, exists := ip[realIP]
//
//			if !exists {
//				ip[realIP] = 1
//			} else {
//				ip[realIP] = count + 1
//
//				if ip[realIP] >= 3 {
//					return c.JSON(http.StatusTooManyRequests, map[string]string{
//						"message": "Слишком много подключений",
//					})
//				}
//			}
//
//			return next(c)
//		}
//	}
//}
