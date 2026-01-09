package handle

import (
	"github.com/Parnishkaspb/curs-abds/internal/service/models"
	"net/http"
	"time"

	"github.com/Parnishkaspb/curs-abds/internal/kafka"
	"github.com/Parnishkaspb/curs-abds/internal/service"
	"github.com/Parnishkaspb/curs-abds/internal/service/frauds"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc    *service.DBService
	frauds *frauds.Frauds
}

func NewHandle(svc *service.DBService, fr *frauds.Frauds) *Handler {
	return &Handler{
		svc:    svc,
		frauds: fr,
	}
}

func (h *Handler) CreateTransaction(c echo.Context) error {
	ctx := c.Request().Context()

	var req kafka.TransactionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Неверный запрос"})
	}

	decision := h.frauds.Evaluate(ctx, req)

	id, err := h.frauds.ApplyDecision(ctx, req, decision, 2)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	status := http.StatusOK
	result := "ACCEPT"
	if decision.Decline {
		status = http.StatusOK
		result = "DECLINE"
	}

	return c.JSON(status, map[string]any{
		"result": result,
		"id":     id,
		"source": 2,
		"reason": decision.Reason,
		"logs":   decision.Logs,
	})
}

func (h *Handler) GetTransactions(c echo.Context) error {
	ctx := c.Request().Context()

	var q models.TransactionsQuery
	if err := c.Bind(&q); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Неверные query-параметры"})
	}

	filter := models.TransactionFilter{
		ID:            q.ID,
		TransactionID: q.TransactionID,
		AccountID:     q.AccountID,
		StatusID:      q.StatusID,
		SourceID:      q.SourceID,
		CountryID:     q.CountryID,
		Merchant:      q.Merchant,
		Accepted:      q.Accepted,
		Limit:         q.Limit,
		Offset:        q.Offset,
	}

	// RFC3339 время -> *time.Time
	if q.CreatedFrom != "" {
		t, err := time.Parse(time.RFC3339, q.CreatedFrom)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "created_from должен быть RFC3339"})
		}
		filter.CreatedFrom = &t
	}
	if q.CreatedTo != "" {
		t, err := time.Parse(time.RFC3339, q.CreatedTo)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "created_to должен быть RFC3339"})
		}
		filter.CreatedTo = &t
	}

	res, err := h.svc.SearchTransactions(ctx, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"total":  res.Total,
		"limit":  res.Limit,
		"offset": res.Offset,
		"items":  res.Items,
	})
}
