package balance_delivery

import (
	"github.com/gofiber/fiber/v2"
	"gofermart/internal/balance"
	"gofermart/pkg/middlewares"
)

func MapBalanceRoutes(group fiber.Router, h balance.Handler, mw *middlewares.MiddlewareManager) {
	group.Get("/balance", mw.AuthMiddleware(), h.GetBalance())
	group.Get("/balance/withdraw", mw.AuthMiddleware(), h.Withdraw())
	group.Get("/withdrawals", mw.AuthMiddleware(), h.GetWithdrawals())
}
