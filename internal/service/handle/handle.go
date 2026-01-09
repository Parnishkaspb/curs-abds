package handle

import (
	"github.com/Parnishkaspb/curs-abds/internal/kafka"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Handle struct {
	db *gorm.DB
}

func NewHandle(db *gorm.DB) *Handle {
	return &Handle{db: db}
}

func (h *Handle) CreateTransaction(c echo.Context) error {
	var req kafka.TransactionRequest

	if err := c.Bind(&req); err != nil {
		log.Printf("Ошибка: %s", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Неверный запрос"})
	}

	//idCountry, exists := Countries[req.Country]
	//if !exists {
	//	log.Printf("Страна не найдена: %s. Автоматическое добавление!", req.Country)
	//	idCountry = maxID() + 1
	//	Countries[req.Country] = idCountry
	//}
	//
	//idCurrency, exists := Currencies[req.Country]
	//if !exists {
	//	log.Fatalf("Валюта не найдена: %s. Автоматическое добавление невозможно!", req.Country)
	//}
	//
	//payloadJSON, err := json.Marshal(req)
	//if err != nil {
	//	log.Fatalf("Ошибка при конвертации данных!")
	//}
	//
	//newTransaction := models.Transaction{
	//	ID:            maxIDTransaction() + 1,
	//	TransactionID: pgtype.UUID{},
	//	AccountID:     req.AccountID,
	//	Amount:        req.Amount,
	//	CurrencyID:    idCurrency,
	//	Merchant:      req.Merchant,
	//	CountryID:     uint64(idCountry),
	//	StatusID:      1,
	//	Payload:       string(payloadJSON),
	//	SourceID:      1,
	//	CreatedAt:     req.CreatedAt,
	//}
	//
	//transactions = append(transactions, newTransaction)
	//
	return c.JSON(http.StatusOK, map[string]string{"message": "Успешное добавление транзакции"})

}
