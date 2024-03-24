package userdelivery

import (
	"github.com/MaxBoych/gofermart/internal/user"
	"github.com/MaxBoych/gofermart/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func MapUserRoutes(group fiber.Router, h user.Handlers, _ *middlewares.MiddlewareManager) {
	group.Post("/register", h.Register())
	group.Post("/login", h.Login())
}
