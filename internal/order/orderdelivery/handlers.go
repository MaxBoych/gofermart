package orderdelivery

import (
	"encoding/json"
	"github.com/MaxBoych/gofermart/internal/order"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type OrderHandler struct {
	useCase order.UseCase
}

func NewOrderHandler(useCase order.UseCase) *OrderHandler {
	return &OrderHandler{
		useCase: useCase,
	}
}

func (h *OrderHandler) UploadNewOrder() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		orderNumber := string(ctx.Body())
		userID := ctx.Locals("user_id").(int64)

		if err := h.useCase.UploadNewOrder(ctx.Context(), orderNumber, userID); err != nil {
			return err
		}

		ctx.Status(http.StatusAccepted)
		return ctx.JSON(fiber.Map{
			"data": "Successfully uploaded",
		})
	}
}

func (h *OrderHandler) GetOrders() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "application/json")
		userID := ctx.Locals("user_id").(int64)

		orders, err := h.useCase.RefreshAndGetOrders(ctx.Context(), userID)
		if err != nil {
			return err
		}

		jsonResp, err := json.Marshal(orders)
		if err != nil {
			return err
		}
		_, err = ctx.Write(jsonResp)

		return err
	}
}
