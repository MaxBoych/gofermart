package user

import (
	"context"
	"gofermart/internal/user/user_models"
)

type UseCase interface {
	Register(ctx context.Context, req user_models.UserRegisterRequest) error
	Login(ctx context.Context, req user_models.UserLoginRequest) (string, error)
}
