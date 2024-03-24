package user_delivery

import (
	"github.com/MaxBoych/gofermart/internal/user"
	"github.com/MaxBoych/gofermart/internal/user/user_models"
	"github.com/MaxBoych/gofermart/pkg/cookie"
	"github.com/MaxBoych/gofermart/pkg/errs"
	"github.com/gofiber/fiber/v2"
	"time"
)

type UserHandler struct {
	useCase user.UseCase
}

func NewUserHandler(useCase user.UseCase) *UserHandler {
	return &UserHandler{
		useCase: useCase,
	}
}

func (h *UserHandler) Register() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		req := user_models.UserRegisterRequest{}
		if err := ctx.BodyParser(&req); err != nil {
			return errs.HttpErrInvalidRequest
		}

		token, err := h.useCase.Register(ctx.Context(), req)
		if err != nil {
			return err
		}

		cookie.SetCookie(ctx, user_models.CookieData{
			Name:    "token",
			Value:   token,
			Expires: time.Now().Add(time.Hour * 72),
			Domain:  ctx.Hostname(),
		})

		ctx.Set("Authorization", token)

		return ctx.JSON(fiber.Map{
			"data": "Successfully registered",
		})
	}
}

func (h *UserHandler) Login() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		req := user_models.UserLoginRequest{}
		if err := ctx.BodyParser(&req); err != nil {
			return errs.HttpErrInvalidRequest
		}

		token, err := h.useCase.Login(ctx.Context(), req)
		if err != nil {
			return err
		}

		cookie.SetCookie(ctx, user_models.CookieData{
			Name:    "token",
			Value:   token,
			Expires: time.Now().Add(time.Hour * 72),
			Domain:  ctx.Hostname(),
		})

		ctx.Set("Authorization", token)

		return ctx.JSON(fiber.Map{
			"data": "Successfully logined",
		})
	}
}
