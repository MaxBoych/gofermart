package balance_usecase

import (
	"context"
	"github.com/MaxBoych/gofermart/internal/balance"
	"github.com/MaxBoych/gofermart/internal/balance/balance_models"
	"github.com/MaxBoych/gofermart/pkg/errs"
	"github.com/MaxBoych/gofermart/pkg/logger"
	"github.com/MaxBoych/gofermart/pkg/utils"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"go.uber.org/zap"
)

type BalanceUC struct {
	balanceRepo balance.Repository
	trManager   *manager.Manager
}

func NewBalanceUC(
	balanceRepo balance.Repository,
	trManager *manager.Manager,
) *BalanceUC {
	return &BalanceUC{
		balanceRepo: balanceRepo,
		trManager:   trManager,
	}
}

func (uc *BalanceUC) GetBalance(ctx context.Context, userId int64) (*balance_models.BalanceResponseData, error) {
	balanceStorage, err := uc.balanceRepo.GetBalance(ctx, userId)
	if err != nil {
		return nil, err
	}

	response := balance_models.BalanceStorageToResponse(*balanceStorage)
	return &response, nil
}

func (uc *BalanceUC) Withdraw(ctx context.Context, req balance_models.WithdrawRequestData) error {
	if !utils.ValidateLuhn(req.Order) {
		return errs.HttpErrOrderIncorrectNumber
	}

	return uc.trManager.Do(ctx, func(ctx context.Context) error {
		currentBalance, err := uc.GetBalance(ctx, req.UserId)
		if err != nil {
			return err
		}
		if currentBalance.Current < req.Sum {
			return errs.HttpErrNotEnoughFunds
		}

		if err := uc.balanceRepo.CreateWithdraw(ctx, req); err != nil {
			return err
		}

		changeData := balance_models.BalanceChangeData{
			Action: "-",
			Sum:    req.Sum,
			UserId: req.UserId,
		}
		logger.Log.Info("withdraw info", zap.Int64("userID", req.UserId), zap.Float64("sum", req.Sum))
		return uc.balanceRepo.UpdateBalance(ctx, changeData)
	})
}

func (uc *BalanceUC) GetWithdrawals(ctx context.Context, userId int64) ([]balance_models.WithdrawResponseData, error) {
	withdrawals, err := uc.balanceRepo.GetWithdrawals(ctx, userId)
	if err != nil {
		return nil, err
	}

	if len(withdrawals) == 0 {
		return nil, errs.HttpErrOrderNoContent
	}

	response := balance_models.WithdrawStorageToResponse(withdrawals)

	return response, nil
}
