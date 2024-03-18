package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gofermart/pkg/cookie"
	"gofermart/pkg/errs"
	"gofermart/pkg/jwt"
	"gofermart/pkg/logger"
)

func (m *MiddlewareManager) AuthMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		key, err := m.tokenRepo.GetSecretKey(ctx.Context())
		if err != nil {
			logger.Log.Error("Error to get secret key", zap.Error(err))
			return errs.HttpErrInternal
		}

		tokenValue, err := cookie.GetCookie(ctx, "token")
		if err != nil {
			logger.Log.Error("Error to get cookie", zap.Error(err))
			return err
		}

		userId, err := jwt.ValidateTokenAndGetUserId(tokenValue, key)
		if err != nil {
			return errs.HttpErrTokenIncorrect
		}

		ctx.Locals("user_id", userId)

		return ctx.Next()
	}
}
