package balance

import (
	"context"
	"github.com/MaxBoych/gofermart/internal/balance/balance_models"
)

type UseCase interface {
	GetBalance(ctx context.Context, userId int64) (*balance_models.BalanceResponseData, error)
	Withdraw(ctx context.Context, req balance_models.WithdrawRequestData) error
	GetWithdrawals(ctx context.Context, userId int64) ([]balance_models.WithdrawResponseData, error)
}
