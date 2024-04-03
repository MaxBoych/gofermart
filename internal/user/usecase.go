package user

import (
	"context"
	"github.com/MaxBoych/gofermart/internal/user/usermodels"
)

type UseCase interface {
	Register(ctx context.Context, req usermodels.UserRegisterRequest) (string, error)
	Login(ctx context.Context, req usermodels.UserLoginRequest) (string, error)
}
