package order

import (
	"context"
	"github.com/MaxBoych/gofermart/internal/order/ordermodels"
)

type UseCase interface {
	UploadNewOrder(ctx context.Context, number string, userID int64) error
	RefreshAndGetOrders(ctx context.Context, userID int64) ([]ordermodels.OrderResponseData, error)
}
