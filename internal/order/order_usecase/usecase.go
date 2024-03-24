package order_usecase

import (
	"context"
	"errors"
	"github.com/MaxBoych/gofermart/internal/accrual_service/client"
	"github.com/MaxBoych/gofermart/internal/balance"
	"github.com/MaxBoych/gofermart/internal/balance/balance_models"
	"github.com/MaxBoych/gofermart/internal/config"
	"github.com/MaxBoych/gofermart/internal/order"
	"github.com/MaxBoych/gofermart/internal/order/order_models"
	"github.com/MaxBoych/gofermart/pkg/errs"
	"github.com/MaxBoych/gofermart/pkg/utils"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5"
	"sync"
)

type OrderUseCase struct {
	orderRepo            order.Repository
	balanceRepo          balance.Repository
	accrualServiceClient *client.AccrualServiceClient
	cfg                  *config.Config
	trManager            *manager.Manager
}

func NewOrderUC(
	orderRepo order.Repository,
	balanceRepo balance.Repository,
	accrualServiceClient *client.AccrualServiceClient,
	cfg *config.Config,
	trManager *manager.Manager,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:            orderRepo,
		balanceRepo:          balanceRepo,
		accrualServiceClient: accrualServiceClient,
		cfg:                  cfg,
		trManager:            trManager,
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

func (uc *OrderUseCase) updateOrder(ctx context.Context, updatedOrder order_models.OrderStorageData) error {
	err := uc.orderRepo.UpdateOrder(ctx, updatedOrder)
	if err != nil {
		return err
	}
	if !updatedOrder.IsProcessed() {
		return nil
	}

	changeData := balance_models.BalanceChangeData{
		Action: "+",
		Sum:    *updatedOrder.Accrual,
		UserId: updatedOrder.UserId,
	}
	return uc.balanceRepo.UpdateBalance(ctx, changeData)
}

func (uc *OrderUseCase) RefreshAndGetOrders(ctx context.Context, userId int64) ([]order_models.OrderResponseData, error) {
	orders, err := uc.orderRepo.GetOrders(ctx, userId)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, errs.HttpErrOrderNoContent
	}

	ordersMutex := struct {
		Orders []order_models.OrderStorageData
		Mutex  sync.Mutex
	}{
		Orders: make([]order_models.OrderStorageData, len(orders)),
	}

	if err = uc.trManager.Do(ctx, func(ctx context.Context) error {
		var wg sync.WaitGroup
		for i, oldOrderData := range orders {
			if oldOrderData.IsFinalStatus() {
				ordersMutex.Orders[i] = oldOrderData
				continue
			}

			wg.Add(1)
			go func(i int, oldOrderData order_models.OrderStorageData) {
				defer wg.Done()

				resp, err := uc.accrualServiceClient.SendRequest(oldOrderData)
				if err != nil {
					ordersMutex.Mutex.Lock()
					ordersMutex.Orders[i] = oldOrderData
					ordersMutex.Mutex.Unlock()
					return
				}

				respData, err := uc.accrualServiceClient.HttpResponseToOrderAccrualResponse(resp)
				if err != nil {
					return
				}

				updatedOrder := order_models.OrderStorageData{
					Number:    oldOrderData.Number,
					Status:    respData.Status,
					Accrual:   &respData.Accrual,
					CreatedAt: oldOrderData.CreatedAt,
				}

				err = uc.updateOrder(ctx, updatedOrder)
				if err != nil {
					return
				}

				ordersMutex.Mutex.Lock()
				ordersMutex.Orders[i] = updatedOrder
				ordersMutex.Mutex.Unlock()
			}(i, oldOrderData)
		}
		wg.Wait()

		return nil
	}); err != nil {
		return nil, err
	}

	response := order_models.OrderStorageToResponse(ordersMutex.Orders)

	return response, nil
}

/*func (uc *OrderUseCase) RefreshAndGetOrders(ctx context.Context, userId int64) ([]order_models.OrderStorageData, error) {
	oldOrdersData, err := uc.orderRepo.GetOrders(ctx, userId, true)
	if err != nil {
		return nil, err
	}

	if len(oldOrdersData) == 0 {
		return
	}
}*/
