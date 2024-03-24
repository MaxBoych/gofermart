package balanceusecase

import (
	"context"
	"github.com/MaxBoych/gofermart/internal/balance"
	"github.com/MaxBoych/gofermart/internal/balance/balancemodels"
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

func (uc *BalanceUC) GetBalance(ctx context.Context, userID int64) (*balancemodels.BalanceResponseData, error) {
	balanceStorage, err := uc.balanceRepo.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := balancemodels.BalanceStorageToResponse(*balanceStorage)
	return &response, nil
}

func (uc *BalanceUC) Withdraw(ctx context.Context, req balancemodels.WithdrawRequestData) error {
	if !utils.ValidateLuhn(req.Order) {
		return errs.HTTPErrOrderIncorrectNumber
	}

	return uc.trManager.Do(ctx, func(ctx context.Context) error {
		currentBalance, err := uc.GetBalance(ctx, req.UserID)
		if err != nil {
			return err
		}
		if currentBalance.Current < req.Sum {
			return errs.HTTPErrNotEnoughFunds
		}

		if err := uc.balanceRepo.CreateWithdraw(ctx, req); err != nil {
			return err
		}

		changeData := balancemodels.BalanceChangeData{
			Action: "-",
			Sum:    req.Sum,
			UserID: req.UserID,
		}
		logger.Log.Info("withdraw info", zap.Int64("userID", req.UserID), zap.Float64("sum", req.Sum))
		return uc.balanceRepo.UpdateBalance(ctx, changeData)
	})
}

func (uc *BalanceUC) GetWithdrawals(ctx context.Context, userId int64) ([]balancemodels.WithdrawResponseData, error) {
	withdrawals, err := uc.balanceRepo.GetWithdrawals(ctx, userId)
	if err != nil {
		return nil, err
	}

	if len(withdrawals) == 0 {
		return nil, errs.HTTPErrOrderNoContent
	}

	response := balancemodels.WithdrawStorageToResponse(withdrawals)

	return response, nil
}
