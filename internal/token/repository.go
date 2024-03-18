package token

import (
	"context"
	"github.com/MaxBoych/gofermart/internal/token/token_models"
)

type Repository interface {
	CreateToken(ctx context.Context, token token_models.TokenStorageData) error
	GetToken(ctx context.Context, userId int64) (*token_models.TokenStorageData, error)
	GetSecretKey(ctx context.Context) (string, error)
}
