package balance_delivery

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"gofermart/internal/balance"
	"gofermart/internal/balance/balance_models"
	"gofermart/pkg/errs"
)

type BalanceHandler struct {
	useCase balance.UseCase
}

func NewBalanceHandler(useCase balance.UseCase) *BalanceHandler {
	return &BalanceHandler{
		useCase: useCase,
	}
}

func (h *BalanceHandler) GetBalance() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId := ctx.Locals("user_id").(int64)

		balanceResp, err := h.useCase.GetBalance(ctx.Context(), userId)
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

		return h.useCase.Withdraw(ctx.Context(), req)
	}
}

func (h *BalanceHandler) GetWithdrawals() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId := ctx.Locals("user_id").(int64)

		withdrawals, err := h.useCase.GetWithdrawals(ctx.Context(), userId)

		jsonResp, err := json.Marshal(withdrawals)
		if err != nil {
			return err
		}
		_, err = ctx.Write(jsonResp)

		return err
	}
}
