package balancedelivery

import (
	"encoding/json"
	"errors"
	"github.com/MaxBoych/gofermart/internal/balance"
	"github.com/MaxBoych/gofermart/internal/balance/balancemodels"
	"github.com/MaxBoych/gofermart/internal/order"
	"github.com/MaxBoych/gofermart/pkg/errs"
	"github.com/gofiber/fiber/v2"
)

type BalanceHandler struct {
	balanceUC balance.UseCase
	orderUC   order.UseCase
}

func NewBalanceHandler(
	balanceUC balance.UseCase,
	orderUC order.UseCase,
) *BalanceHandler {
	return &BalanceHandler{
		balanceUC: balanceUC,
		orderUC:   orderUC,
	}
}

func (h *BalanceHandler) GetBalance() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "application/json")
		userID, ok := ctx.Locals("user_id").(int64)
		if !ok {
			return errs.HTTPErrUnauthorized
		}

		c := ctx.Context()

		_, err := h.orderUC.RefreshAndGetOrders(c, userID)
		if err != nil && !errors.Is(err, errs.HTTPErrOrderNoContent) {
			return err
		}

		balanceResp, err := h.balanceUC.GetBalance(c, userID)
		if err != nil {
			return err
		}

		jsonResp, err := json.Marshal(balanceResp)
		if err != nil {
			return err
		}
		_, err = ctx.Write(jsonResp)

		return err
	}
}

func (h *BalanceHandler) Withdraw() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userID, ok := ctx.Locals("user_id").(int64)
		if !ok {
			return errs.HTTPErrUnauthorized
		}

		req := balancemodels.WithdrawRequestData{}
		if err := ctx.BodyParser(&req); err != nil {
			return errs.HTTPErrInvalidRequest
		}
		req.UserID = userID

		c := ctx.Context()

		_, err := h.orderUC.RefreshAndGetOrders(c, userID)
		if err != nil && !errors.Is(err, errs.HTTPErrOrderNoContent) {
			return err
		}

		return h.balanceUC.Withdraw(c, req)
	}
}

func (h *BalanceHandler) GetWithdrawals() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "application/json")
		userID, ok := ctx.Locals("user_id").(int64)
		if !ok {
			return errs.HTTPErrUnauthorized
		}

		withdrawals, err := h.balanceUC.GetWithdrawals(ctx.Context(), userID)
		if err != nil {
			return err
		}

		jsonResp, err := json.Marshal(withdrawals)
		if err != nil {
			return err
		}
		_, err = ctx.Write(jsonResp)

		return err
	}
}
