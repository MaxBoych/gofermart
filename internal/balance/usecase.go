package balance

import (
	"context"
	"github.com/MaxBoych/gofermart/internal/balance/balancemodels"
)

type UseCase interface {
	GetBalance(ctx context.Context, userID int64) (*balancemodels.BalanceResponseData, error)
	Withdraw(ctx context.Context, req balancemodels.WithdrawRequestData) error
	GetWithdrawals(ctx context.Context, userID int64) ([]balancemodels.WithdrawResponseData, error)
}
