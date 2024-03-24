package tokenrepository

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/MaxBoych/gofermart/internal/token/tokenmodels"
	"github.com/MaxBoych/gofermart/pkg/logger"
	"github.com/MaxBoych/gofermart/pkg/sqlqueries"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type TokenRepo struct {
	db       *pgxpool.Pool
	txGetter *trmpgx.CtxGetter
}

func NewTokenRepo(db *pgxpool.Pool, txGetter *trmpgx.CtxGetter) *TokenRepo {
	return &TokenRepo{
		db:       db,
		txGetter: txGetter,
	}
}

func (r *TokenRepo) GetSecretKey(ctx context.Context) (string, error) {
	secretKey := tokenmodels.SecretKeyStorageData{}
	query, args, err := sq.Select(sqlqueries.SelectSecretKey...).
		From(sqlqueries.SecretKeyTableName).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Log.Error("Error to build sql SELECT query", zap.Error(err))
		return "", err
	}

	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	err = tr.QueryRow(ctx, query, args...).Scan(&secretKey.Value, &secretKey.CreatedAt, &secretKey.UpdatedAt)
	if err != nil {
		logger.Log.Error("Error while scanning, sql SELECT query", zap.Error(err))
		return "", err
	}

	return secretKey.Value, nil
}

func (r *TokenRepo) GetToken(ctx context.Context, userID int64) (*tokenmodels.TokenStorageData, error) {
	tokenData := tokenmodels.TokenStorageData{}
	query, args, err := sq.Select(sqlqueries.SelectToken...).
		From(sqlqueries.TokenTableName).
		Where(sq.Eq{sqlqueries.UserIDColumnName: userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Log.Error("Error to build sql SELECT query", zap.Error(err))
		return nil, err
	}

	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	err = tr.QueryRow(ctx, query, args...).Scan(&tokenData.TokenID, &tokenData.UserID, &tokenData.Value, &tokenData.CreatedAt, &tokenData.UpdatedAt)
	if err != nil {
		logger.Log.Error("Error while scanning, sql SELECT query", zap.Error(err))
		return nil, err
	}

	return &tokenData, nil
}

func (r *TokenRepo) CreateToken(ctx context.Context, token tokenmodels.TokenStorageData) error {
	query, args, err := sq.Insert(sqlqueries.TokenTableName).
		Columns(sqlqueries.InsertToken...).
		Values(
			token.UserID,
			token.Value,
			time.Now(),
			time.Now(),
		).
		Suffix(fmt.Sprintf("ON CONFLICT (%[1]s) DO UPDATE SET %[2]s = EXCLUDED.%[2]s, %[3]s = EXCLUDED.%[3]s",
			sqlqueries.UserIDColumnName,      // 1
			sqlqueries.TokenValueColumnName,  // 2
			sqlqueries.UpdatedAtColumnName)). // 3
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Log.Error("Error to build sql INSERT query", zap.Error(err))
		return err
	}

	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	_, err = tr.Exec(ctx, query, args...)
	if err != nil {
		logger.Log.Error("Error to execute sql INSERT query", zap.Error(err))
		return err
	}

	return nil
}
