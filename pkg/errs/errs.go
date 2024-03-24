package errs

import "errors"

var (
	ErrInvalidRequest            = errors.New("bad request was provided")
	ErrUserAlreadyExists         = errors.New("user already exists")
	ErrUserIncorrectLogin        = errors.New("incorrect login or password")
	ErrOrderIncorrectNumber      = errors.New("incorrect order number")
	ErrCookieIsEmpty             = errors.New("cookie is empty")
	ErrTokenIncorrect            = errors.New("incorrect token")
	ErrInternal                  = errors.New("internal server error")
	ErrOrderDuplicateSameUser    = errors.New("order was already uploaded by this user")
	ErrOrderDuplicateAnotherUser = errors.New("order was already uploaded by another user")
	ErrOrderNoContent            = errors.New("this user has no such data")
	ErrNotEnoughFunds            = errors.New("not enough funds on balance")
	ErrTooManyRequests           = errors.New("too many requests")
	ErrConnectionRefused         = errors.New("connection refused")
)
