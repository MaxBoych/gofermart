package balance

import (
	"context"
	"gofermart/internal/balance/balance_models"
)

type Repository interface {
	GetBalance(ctx context.Context, userId int64) (*balance_models.BalanceStorageData, error)
	CreateBalance(ctx context.Context, userId int64) error
	UpdateBalance(ctx context.Context, req balance_models.WithdrawRequestData) error
	CreateWithdraw(ctx context.Context, req balance_models.WithdrawRequestData) error
	GetWithdrawals(ctx context.Context, userId int64) ([]balance_models.WithdrawStorageData, error)
}
