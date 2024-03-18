package middlewares

import "gofermart/internal/token"

type MiddlewareManager struct {
	tokenRepo token.Repository
}

func NewMiddlewareManager(tokenRepo token.Repository) *MiddlewareManager {
	return &MiddlewareManager{
		tokenRepo: tokenRepo,
	}
}
