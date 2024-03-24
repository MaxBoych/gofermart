package order

import (
	"context"
	"github.com/MaxBoych/gofermart/internal/order/order_models"
)

type UseCase interface {
	UploadNewOrder(ctx context.Context, number string, userId int64) error
	RefreshAndGetOrders(ctx context.Context, userId int64) ([]order_models.OrderResponseData, error)
}
