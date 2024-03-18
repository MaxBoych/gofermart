package main

import (
	"context"
	"fmt"
	"github.com/MaxBoych/gofermart/internal/config"
	"github.com/MaxBoych/gofermart/internal/server"
	database "github.com/MaxBoych/gofermart/pkg/db"
	"github.com/MaxBoych/gofermart/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	if err := logger.Initialize("INFO"); err != nil {
		fmt.Printf("Error to init logger: %v\n", err)
		return
	}

	cfg := config.NewConfig()
	cfg.ParseConfig()

	db := database.NewDB()
	if err := db.Connect(ctx, cfg.DatabaseDSN); err != nil {
		logger.Log.Error("Error to connect database", zap.Error(err))
		return
	}
	if err := db.Init(ctx); err != nil {
		logger.Log.Error("Error to init database", zap.Error(err))
		return
	}

	s := server.NewServer(cfg, db)
	if err := s.Run(); err != nil {
		logger.Log.Error("Error to start server", zap.Error(err))
		return
	}
}
