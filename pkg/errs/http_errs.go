package errs

import "github.com/gofiber/fiber/v2"

var (
	HttpErrInvalidRequest            = fiber.NewError(fiber.StatusBadRequest, ErrInvalidRequest.Error())
	HttpErrUserAlreadyExists         = fiber.NewError(fiber.StatusConflict, ErrUserAlreadyExists.Error())
	HttpErrUserIncorrectLogin        = fiber.NewError(fiber.StatusUnauthorized, ErrUserIncorrectLogin.Error())
	HttpErrOrderIncorrectNumber      = fiber.NewError(fiber.StatusUnprocessableEntity, ErrOrderIncorrectNumber.Error())
	HttpErrCookieIsEmpty             = fiber.NewError(fiber.StatusUnauthorized, ErrCookieIsEmpty.Error())
	HttpErrTokenIncorrect            = fiber.NewError(fiber.StatusUnauthorized, ErrTokenIncorrect.Error())
	HttpErrInternal                  = fiber.NewError(fiber.StatusInternalServerError, ErrInternal.Error())
	HttpErrOrderDuplicateSameUser    = fiber.NewError(fiber.StatusOK, ErrOrderDuplicateSameUser.Error())
	HttpErrOrderDuplicateAnotherUser = fiber.NewError(fiber.StatusConflict, ErrOrderDuplicateAnotherUser.Error())
	HttpErrOrderNoContent            = fiber.NewError(fiber.StatusNoContent, ErrOrderNoContent.Error())
	HttpErrNotEnoughFunds            = fiber.NewError(fiber.StatusPaymentRequired, ErrNotEnoughFunds.Error())
	HttpErrTooManyRequests           = fiber.NewError(fiber.StatusTooManyRequests, fiber.ErrTooManyRequests.Error())
	HttpErrConnectionRefused         = fiber.NewError(fiber.StatusNotAcceptable, fiber.ErrTooManyRequests.Error())
)
