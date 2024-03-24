package accrual_service_models

import (
	"github.com/MaxBoych/gofermart/internal/order/order_models"
	"net/http"
)

type AccrualRequestWithResponse struct {
	Request  *http.Request
	Response chan *http.Response
}

type AccrualOrderResponseData struct {
	Number  string                   `json:"number"`
	Status  order_models.OrderStatus `json:"status"`
	Accrual float64                  `json:"accrual"`
}
