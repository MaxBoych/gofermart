package server

import (
	"github.com/MaxBoych/gofermart/internal/balance/balance_delivery"
	"github.com/MaxBoych/gofermart/internal/balance/balance_repository"
	"github.com/MaxBoych/gofermart/internal/balance/balance_usecase"
	"github.com/MaxBoych/gofermart/internal/order/order_delivery"
	"github.com/MaxBoych/gofermart/internal/order/order_repository"
	"github.com/MaxBoych/gofermart/internal/order/order_usecase"
	"github.com/MaxBoych/gofermart/internal/token/token_repository"
	"github.com/MaxBoych/gofermart/internal/user/user_delivery"
	"github.com/MaxBoych/gofermart/internal/user/user_repository"
	"github.com/MaxBoych/gofermart/internal/user/user_usecase"
	"github.com/MaxBoych/gofermart/pkg/middlewares"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

func (s *Server) MapHandlers() {
	trManager := manager.Must(trmpgx.NewDefaultFactory(s.db.Pool))
	getter := trmpgx.DefaultCtxGetter

	tokenRepo := token_repository.NewTokenRepo(s.db.Pool, getter)
	balanceRepo := balance_repository.NewBalanceRepo(s.db.Pool, getter)
	mw := middlewares.NewMiddlewareManager(tokenRepo)

	userRepo := user_repository.NewUserRepo(s.db.Pool, getter)
	userUC := user_usecase.NewUserUC(userRepo, tokenRepo, balanceRepo, trManager)
	userHandler := user_delivery.NewUserHandler(userUC)
	userGroup := s.fb.Group("api/user")
	user_delivery.MapUserRoutes(userGroup, userHandler, mw)

	orderRepo := order_repository.NewOrderRepo(s.db.Pool, getter)
	orderUC := order_usecase.NewOrderUC(orderRepo, balanceRepo, s.accrualServiceClient, s.cfg, trManager)
	orderHandler := order_delivery.NewOrderHandler(orderUC)
	orderGroup := s.fb.Group("api/user/orders")
	order_delivery.MapOrderRoutes(orderGroup, orderHandler, mw)

	balanceUC := balance_usecase.NewBalanceUC(balanceRepo, trManager)
	balanceHandler := balance_delivery.NewBalanceHandler(balanceUC, orderUC)
	balanceGroup := s.fb.Group("api/user")
	balance_delivery.MapBalanceRoutes(balanceGroup, balanceHandler, mw)
}
