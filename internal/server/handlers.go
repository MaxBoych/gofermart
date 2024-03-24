package server

import (
	"github.com/MaxBoych/gofermart/internal/balance/balancedelivery"
	"github.com/MaxBoych/gofermart/internal/balance/balancerepository"
	"github.com/MaxBoych/gofermart/internal/balance/balanceusecase"
	"github.com/MaxBoych/gofermart/internal/order/orderdelivery"
	"github.com/MaxBoych/gofermart/internal/order/orderrepository"
	"github.com/MaxBoych/gofermart/internal/order/orderusecase"
	"github.com/MaxBoych/gofermart/internal/token/tokenrepository"
	"github.com/MaxBoych/gofermart/internal/user/userdelivery"
	"github.com/MaxBoych/gofermart/internal/user/userrepository"
	"github.com/MaxBoych/gofermart/internal/user/userusecase"
	"github.com/MaxBoych/gofermart/pkg/middlewares"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

func (s *Server) MapHandlers() {
	trManager := manager.Must(trmpgx.NewDefaultFactory(s.db.Pool))
	getter := trmpgx.DefaultCtxGetter

	tokenRepo := tokenrepository.NewTokenRepo(s.db.Pool, getter)
	balanceRepo := balancerepository.NewBalanceRepo(s.db.Pool, getter)
	mw := middlewares.NewMiddlewareManager(tokenRepo)

	userRepo := userrepository.NewUserRepo(s.db.Pool, getter)
	userUC := userusecase.NewUserUC(userRepo, tokenRepo, balanceRepo, trManager)
	userHandler := userdelivery.NewUserHandler(userUC)
	userGroup := s.fb.Group("api/user")
	userdelivery.MapUserRoutes(userGroup, userHandler, mw)

	orderRepo := orderrepository.NewOrderRepo(s.db.Pool, getter)
	orderUC := orderusecase.NewOrderUC(orderRepo, balanceRepo, s.accrualServiceClient, s.cfg, trManager)
	orderHandler := orderdelivery.NewOrderHandler(orderUC)
	orderGroup := s.fb.Group("api/user/orders")
	orderdelivery.MapOrderRoutes(orderGroup, orderHandler, mw)

	balanceUC := balanceusecase.NewBalanceUC(balanceRepo, trManager)
	balanceHandler := balancedelivery.NewBalanceHandler(balanceUC, orderUC)
	balanceGroup := s.fb.Group("api/user")
	balancedelivery.MapBalanceRoutes(balanceGroup, balanceHandler, mw)
}
