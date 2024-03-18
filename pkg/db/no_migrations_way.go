package db

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/MaxBoych/gofermart/internal/store/sql_queries"
	"github.com/MaxBoych/gofermart/pkg/jwt"
	"github.com/MaxBoych/gofermart/pkg/logger"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"go.uber.org/zap"
)

func (db *DB) Init(ctx context.Context) error {
	trManager := manager.Must(trmpgx.NewDefaultFactory(db.Pool))
	if err := trManager.Do(ctx, func(ctx context.Context) error {

		tr := trmpgx.DefaultCtxGetter.DefaultTrOrDB(ctx, db.Pool)

		// Create secret key table
		_, err := tr.Exec(ctx, sql_queries.SecretKeyTableSQL)
		if err != nil {
			logger.Log.Error("Error to create secret key table", zap.Error(err))
			return err
		}
		secretKey, err := jwt.GenerateSecretKey()
		if err != nil {
			logger.Log.Error("Error to generate secret key for jwt", zap.Error(err))
			return err
		}
		query, args, err := sq.Insert(sql_queries.SecretKeyTableName).
			Columns(sql_queries.InsertSecretKey...).
			Values(1, secretKey, sq.Expr("NOW()"), sq.Expr("NOW()")).
			PlaceholderFormat(sq.Dollar).
			Suffix(fmt.Sprintf("ON CONFLICT (%[1]s) DO NOTHING", sql_queries.SecretKeyConstantColumnName)).
			ToSql()
		if err != nil {
			logger.Log.Error("Error to build INSERT query", zap.Error(err))
			return err
		}
		_, err = tr.Exec(ctx, query, args...)
		if err != nil {
			logger.Log.Error("Error to execute INSERT query", zap.Error(err))
			return err
		}

		// Create jwt table
		_, err = tr.Exec(ctx, sql_queries.JwtTableSQL)
		if err != nil {
			logger.Log.Error("Error to create jwt table", zap.Error(err))
			return err
		}

		// Create user table
		_, err = tr.Exec(ctx, sql_queries.UserTableSQL)
		if err != nil {
			logger.Log.Error("Error to create user table", zap.Error(err))
			return err
		}

		// Create order table
		_, err = tr.Exec(ctx, sql_queries.OrderTableSQL)
		if err != nil {
			logger.Log.Error("Error to create order table", zap.Error(err))
			return err
		}

		// Create balance table
		_, err = tr.Exec(ctx, sql_queries.BalanceTableSQL)
		if err != nil {
			logger.Log.Error("Error to create balance table", zap.Error(err))
			return err
		}

		// Create withdraw table
		_, err = tr.Exec(ctx, sql_queries.WithdrawTableSQL)
		if err != nil {
			logger.Log.Error("Error to create withdraw table", zap.Error(err))
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	logger.Log.Info("DB tables created successfully")
	return nil
}
