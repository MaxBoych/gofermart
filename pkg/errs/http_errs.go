package errs

import "github.com/gofiber/fiber/v2"

var (
	HTTPErrInvalidRequest            = fiber.NewError(fiber.StatusBadRequest, ErrInvalidRequest.Error())
	HTTPErrUserAlreadyExists         = fiber.NewError(fiber.StatusConflict, ErrUserAlreadyExists.Error())
	HTTPErrUserIncorrectLogin        = fiber.NewError(fiber.StatusUnauthorized, ErrUserIncorrectLogin.Error())
	HTTPErrOrderIncorrectNumber      = fiber.NewError(fiber.StatusUnprocessableEntity, ErrOrderIncorrectNumber.Error())
	HTTPErrCookieIsEmpty             = fiber.NewError(fiber.StatusUnauthorized, ErrCookieIsEmpty.Error())
	HTTPErrTokenIncorrect            = fiber.NewError(fiber.StatusUnauthorized, ErrTokenIncorrect.Error())
	HTTPErrInternal                  = fiber.NewError(fiber.StatusInternalServerError, ErrInternal.Error())
	HTTPErrOrderDuplicateSameUser    = fiber.NewError(fiber.StatusOK, ErrOrderDuplicateSameUser.Error())
	HTTPErrOrderDuplicateAnotherUser = fiber.NewError(fiber.StatusConflict, ErrOrderDuplicateAnotherUser.Error())
	HTTPErrOrderNoContent            = fiber.NewError(fiber.StatusNoContent, ErrOrderNoContent.Error())
	HTTPErrNotEnoughFunds            = fiber.NewError(fiber.StatusPaymentRequired, ErrNotEnoughFunds.Error())
	HTTPErrTooManyRequests           = fiber.NewError(fiber.StatusTooManyRequests, fiber.ErrTooManyRequests.Error())
	HTTPErrConnectionRefused         = fiber.NewError(fiber.StatusNotAcceptable, fiber.ErrTooManyRequests.Error())
)
