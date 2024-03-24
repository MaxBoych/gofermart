package order

import (
	"context"
	"github.com/MaxBoych/gofermart/internal/order/ordermodels"
)

type Repository interface {
	CreateOrder(ctx context.Context, data ordermodels.OrderStorageData) error
	GetOrder(ctx context.Context, number string) (*ordermodels.OrderStorageData, error)
	GetOrders(ctx context.Context, userID int64) ([]ordermodels.OrderStorageData, error)
	UpdateOrder(ctx context.Context, updatedOrder ordermodels.OrderStorageData) error
}
