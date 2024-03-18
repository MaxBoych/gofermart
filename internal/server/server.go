package server

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gofermart/internal/config"
	"gofermart/pkg/db"
	"gofermart/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	cfg *config.Config
	db  *db.DB
	fb  *fiber.App
}

func NewServer(cfg *config.Config, db *db.DB) *Server {
	return &Server{
		cfg: cfg,
		db:  db,
		fb:  fiber.New(),
	}
}

func (s *Server) Run() error {
	s.MapHandlers()

	go func() {
		addr := s.cfg.RunAddr
		logger.Log.Info("Server is started", zap.String("address", addr))

		err := s.fb.Listen(addr)
		if err != nil {
			logger.Log.Fatal(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit

	err := s.fb.Shutdown()
	if err != nil {
		logger.Log.Error("Error to shutdown fiber", zap.Error(err))
	} else {
		logger.Log.Info("Fiber closed properly")
	}

	return nil
}
