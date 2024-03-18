package user_delivery

import (
	"github.com/gofiber/fiber/v2"
	"gofermart/internal/user"
	"gofermart/pkg/middlewares"
)

func MapUserRoutes(group fiber.Router, h user.Handlers, _ *middlewares.MiddlewareManager) {
	group.Post("/register", h.Register())
	group.Post("/login", h.Login())
}
