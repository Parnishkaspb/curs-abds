package handle

import (
	"net/http"

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
