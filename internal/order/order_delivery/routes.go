package order_delivery

import (
	"github.com/gofiber/fiber/v2"
	"gofermart/internal/order"
	"gofermart/pkg/middlewares"
)

func MapOrderRoutes(group fiber.Router, h order.Handlers, mw *middlewares.MiddlewareManager) {
	group.Post("/", mw.AuthMiddleware(), h.UploadNewOrder())
	group.Get("/", mw.AuthMiddleware(), h.GetOrders())
}
