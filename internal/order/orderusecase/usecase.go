package orderusecase

import (
	"context"
	"errors"
	"github.com/MaxBoych/gofermart/internal/accrual_service/client"
	"github.com/MaxBoych/gofermart/internal/balance"
	"github.com/MaxBoych/gofermart/internal/balance/balancemodels"
	"github.com/MaxBoych/gofermart/internal/config"
	"github.com/MaxBoych/gofermart/internal/order"
	"github.com/MaxBoych/gofermart/internal/order/ordermodels"
	"github.com/MaxBoych/gofermart/pkg/errs"
	"github.com/MaxBoych/gofermart/pkg/logger"
	"github.com/MaxBoych/gofermart/pkg/utils"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
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

func (uc *OrderUseCase) UploadNewOrder(ctx context.Context, number string, userID int64) error {
	if !utils.ValidateLuhn(number) {
		return errs.HTTPErrOrderIncorrectNumber
	}

	newOrder := ordermodels.OrderStorageData{
		Number:  number,
		UserID:  userID,
		Status:  ordermodels.OrderStatusNew,
		Accrual: nil,
	}

	err := uc.trManager.Do(ctx, func(ctx context.Context) error {
		existedOrder, err := uc.orderRepo.GetOrder(ctx, number)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		if existedOrder != nil {
			if newOrder.UserID != existedOrder.UserID {
				return errs.HTTPErrOrderDuplicateAnotherUser
			}

			return errs.HTTPErrOrderDuplicateSameUser
		}

		return uc.orderRepo.CreateOrder(ctx, newOrder)
	})

	return err
}

func (uc *OrderUseCase) updateOrder(ctx context.Context, updatedOrder ordermodels.OrderStorageData) error {
	err := uc.orderRepo.UpdateOrder(ctx, updatedOrder)
	if err != nil {
		return err
	}
	if !updatedOrder.IsProcessed() {
		return nil
	}

	logger.Log.Info("updatedOrder info", zap.Int64("userID", updatedOrder.UserID), zap.Float64("accrual", *updatedOrder.Accrual))
	changeData := balancemodels.BalanceChangeData{
		Action: "+",
		Sum:    *updatedOrder.Accrual,
		UserID: updatedOrder.UserID,
	}
	return uc.balanceRepo.UpdateBalance(ctx, changeData)
}

func (uc *OrderUseCase) RefreshAndGetOrders(ctx context.Context, userID int64) ([]ordermodels.OrderResponseData, error) {
	orders, err := uc.orderRepo.GetOrders(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, errs.HTTPErrOrderNoContent
	}

	ordersMutex := struct {
		Orders []ordermodels.OrderStorageData
		Mutex  sync.Mutex
	}{
		Orders: make([]ordermodels.OrderStorageData, len(orders)),
	}

	if err = uc.trManager.Do(ctx, func(ctx context.Context) error {
		var wg sync.WaitGroup
		for i, oldOrderData := range orders {
			if oldOrderData.IsFinalStatus() {
				ordersMutex.Orders[i] = oldOrderData
				continue
			}

			wg.Add(1)
			go func(i int, oldOrderData ordermodels.OrderStorageData) {
				defer wg.Done()

				resp, err := uc.accrualServiceClient.SendRequest(oldOrderData)
				if err != nil {
					ordersMutex.Mutex.Lock()
					ordersMutex.Orders[i] = oldOrderData
					ordersMutex.Mutex.Unlock()
					return
				}

				respData, err := uc.accrualServiceClient.HTTPResponseToOrderAccrualResponse(resp)
				if err != nil {
					return
				}

				updatedOrder := ordermodels.OrderStorageData{
					Number:    oldOrderData.Number,
					Status:    respData.Status,
					Accrual:   &respData.Accrual,
					UserID:    userID,
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

	response := ordermodels.OrderStorageToResponse(ordersMutex.Orders)

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
