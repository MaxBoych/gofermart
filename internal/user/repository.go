package user

import (
	"context"
	"github.com/MaxBoych/gofermart/internal/user/user_models"
)

type Repository interface {
	CreateUser(ctx context.Context, req user_models.UserStorageData) (int64, error)
	GetUserByLogin(ctx context.Context, login string) (*user_models.UserStorageData, error)
}
