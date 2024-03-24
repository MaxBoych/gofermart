package db

import (
	"context"
	"github.com/MaxBoych/gofermart/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB() *DB {
	return &DB{}
}

func (db *DB) Connect(ctx context.Context, dsn string) error {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}

	logger.Log.Info("Connecting to database", zap.String("address", dsn))
	db.Pool = pool

	return nil
}
