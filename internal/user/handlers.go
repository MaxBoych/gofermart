package user

import "github.com/gofiber/fiber/v2"

type Handlers interface {
	Register() fiber.Handler
	Login() fiber.Handler
}
