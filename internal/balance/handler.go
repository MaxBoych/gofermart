package balance

import "github.com/gofiber/fiber/v2"

type Handler interface {
	GetBalance() fiber.Handler
	Withdraw() fiber.Handler
	GetWithdrawals() fiber.Handler
}
