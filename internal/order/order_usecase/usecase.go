package order_usecase

import (
	"context"
	"errors"
	"github.com/MaxBoych/gofermart/internal/order"
	"github.com/MaxBoych/gofermart/internal/order/order_models"
	"github.com/MaxBoych/gofermart/pkg/errs"
	"github.com/MaxBoych/gofermart/pkg/utils"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5"
)

type OrderUseCase struct {
	orderRepo order.Repository
	trManager *manager.Manager
}

func NewOrderUC(
	orderRepo order.Repository,
	trManager *manager.Manager,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo: orderRepo,
		trManager: trManager,
	}
}

func (uc *OrderUseCase) UploadNewOrder(ctx context.Context, number string, userId int64) error {
	if !utils.ValidateLuhn(number) {
		return errs.HttpErrOrderIncorrectNumber
	}

	newOrder := order_models.OrderStorageData{
		Number:  number,
		UserId:  userId,
		Status:  order_models.OrderStatusNew,
		Accrual: nil,
	}

	err := uc.trManager.Do(ctx, func(ctx context.Context) error {
		existedOrder, err := uc.orderRepo.GetOrder(ctx, number)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		if existedOrder != nil {
			if newOrder.UserId != existedOrder.UserId {
				return errs.HttpErrOrderDuplicateAnotherUser
			} else {
				return errs.HttpErrOrderDuplicateSameUser
			}
		}

		return uc.orderRepo.CreateOrder(ctx, newOrder)
	})

	return err
}

func (uc *OrderUseCase) GetOrders(ctx context.Context, userId int64) ([]order_models.OrderResponseData, error) {
	orders, err := uc.orderRepo.GetOrders(ctx, userId)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, errs.HttpErrOrderNoContent
	}

	response := order_models.OrderStorageToResponse(orders)

	return response, nil
}
