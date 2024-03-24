package token

import (
	"context"
	"github.com/MaxBoych/gofermart/internal/token/tokenmodels"
)

type Repository interface {
	CreateToken(ctx context.Context, token tokenmodels.TokenStorageData) error
	GetToken(ctx context.Context, userID int64) (*tokenmodels.TokenStorageData, error)
	GetSecretKey(ctx context.Context) (string, error)
}
