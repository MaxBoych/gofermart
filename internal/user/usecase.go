package user

import (
	"context"
	"github.com/MaxBoych/gofermart/internal/user/user_models"
)

type UseCase interface {
	Register(ctx context.Context, req user_models.UserRegisterRequest) (string, error)
	Login(ctx context.Context, req user_models.UserLoginRequest) (string, error)
}
