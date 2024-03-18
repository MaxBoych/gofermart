package order

import (
	"context"
	"gofermart/internal/order/order_models"
)

type Repository interface {
	CreateOrder(ctx context.Context, data order_models.OrderStorageData) error
	GetOrder(ctx context.Context, number string) (*order_models.OrderStorageData, error)
	GetOrders(ctx context.Context, userId int64) ([]order_models.OrderStorageData, error)
}
