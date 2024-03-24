package balance_delivery

import (
	"encoding/json"
	"errors"
	"github.com/MaxBoych/gofermart/internal/balance"
	"github.com/MaxBoych/gofermart/internal/balance/balance_models"
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
		userId := ctx.Locals("user_id").(int64)

		c := ctx.Context()

		_, err := h.orderUC.RefreshAndGetOrders(c, userId)
		if err != nil && !errors.Is(err, errs.HttpErrOrderNoContent) {
			return err
		}

		balanceResp, err := h.balanceUC.GetBalance(c, userId)
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
		userId := ctx.Locals("user_id").(int64)

		req := balance_models.WithdrawRequestData{}
		if err := ctx.BodyParser(&req); err != nil {
			return errs.HttpErrInvalidRequest
		}
		req.UserId = userId

		c := ctx.Context()

		_, err := h.orderUC.RefreshAndGetOrders(c, userId)
		if err != nil && !errors.Is(err, errs.HttpErrOrderNoContent) {
			return err
		}

		return h.balanceUC.Withdraw(c, req)
	}
}

func (h *BalanceHandler) GetWithdrawals() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "application/json")
		userId := ctx.Locals("user_id").(int64)

		withdrawals, err := h.balanceUC.GetWithdrawals(ctx.Context(), userId)
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
