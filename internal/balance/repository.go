package balance

import (
	"context"
	"github.com/MaxBoych/gofermart/internal/balance/balancemodels"
)

type Repository interface {
	GetBalance(ctx context.Context, userID int64) (*balancemodels.BalanceStorageData, error)
	CreateBalance(ctx context.Context, userID int64) error
	UpdateBalance(ctx context.Context, req balancemodels.BalanceChangeData) error
	CreateWithdraw(ctx context.Context, req balancemodels.WithdrawRequestData) error
	GetWithdrawals(ctx context.Context, userID int64) ([]balancemodels.WithdrawStorageData, error)
}
