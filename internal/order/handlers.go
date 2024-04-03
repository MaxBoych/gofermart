package order

import "github.com/gofiber/fiber/v2"

type Handlers interface {
	UploadNewOrder() fiber.Handler
	GetOrders() fiber.Handler
}
