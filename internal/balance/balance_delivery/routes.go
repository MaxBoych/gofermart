package balance_delivery

import (
	"github.com/MaxBoych/gofermart/internal/balance"
	"github.com/MaxBoych/gofermart/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func MapBalanceRoutes(group fiber.Router, h balance.Handler, mw *middlewares.MiddlewareManager) {
	group.Get("/balance", mw.AuthMiddleware(), h.GetBalance())
	group.Post("/balance/withdraw", mw.AuthMiddleware(), h.Withdraw())
	group.Get("/withdrawals", mw.AuthMiddleware(), h.GetWithdrawals())
}
