package accrualservicemodels

import (
	"github.com/MaxBoych/gofermart/internal/order/ordermodels"
	"net/http"
)

type AccrualRequestWithResponse struct {
	Request  *http.Request
	Response chan *http.Response
}

type AccrualOrderResponseData struct {
	Number  string                  `json:"number"`
	Status  ordermodels.OrderStatus `json:"status"`
	Accrual float64                 `json:"accrual"`
}
