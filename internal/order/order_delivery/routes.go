package order_delivery

import (
	"github.com/MaxBoych/gofermart/internal/order"
	"github.com/MaxBoych/gofermart/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func MapOrderRoutes(group fiber.Router, h order.Handlers, mw *middlewares.MiddlewareManager) {
	group.Post("/", mw.AuthMiddleware(), h.UploadNewOrder())
	group.Get("/", mw.AuthMiddleware(), h.GetOrders())
}
