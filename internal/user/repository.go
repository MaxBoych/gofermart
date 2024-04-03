package user

import (
	"context"
	"github.com/MaxBoych/gofermart/internal/user/usermodels"
)

type Repository interface {
	CreateUser(ctx context.Context, req usermodels.UserStorageData) (int64, error)
	GetUserByLogin(ctx context.Context, login string) (*usermodels.UserStorageData, error)
}
