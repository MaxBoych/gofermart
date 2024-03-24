package token_repository

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/MaxBoych/gofermart/internal/token/token_models"
	"github.com/MaxBoych/gofermart/pkg/logger"
	"github.com/MaxBoych/gofermart/pkg/sql_queries"
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
	secretKey := token_models.SecretKeyStorageData{}
	query, args, err := sq.Select(sql_queries.SelectSecretKey...).
		From(sql_queries.SecretKeyTableName).
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

func (r *TokenRepo) GetToken(ctx context.Context, userId int64) (*token_models.TokenStorageData, error) {
	tokenData := token_models.TokenStorageData{}
	query, args, err := sq.Select(sql_queries.SelectToken...).
		From(sql_queries.TokenTableName).
		Where(sq.Eq{sql_queries.UserIdColumnName: userId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Log.Error("Error to build sql SELECT query", zap.Error(err))
		return nil, err
	}

	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	err = tr.QueryRow(ctx, query, args...).Scan(&tokenData.TokenId, &tokenData.UserId, &tokenData.Value, &tokenData.CreatedAt, &tokenData.UpdatedAt)
	if err != nil {
		logger.Log.Error("Error while scanning, sql SELECT query", zap.Error(err))
		return nil, err
	}

	return &tokenData, nil
}

func (r *TokenRepo) CreateToken(ctx context.Context, token token_models.TokenStorageData) error {
	query, args, err := sq.Insert(sql_queries.TokenTableName).
		Columns(sql_queries.InsertToken...).
		Values(
			token.UserId,
			token.Value,
			time.Now(),
			time.Now(),
		).
		Suffix(fmt.Sprintf("ON CONFLICT (%[1]s) DO UPDATE SET %[2]s = EXCLUDED.%[2]s, %[3]s = EXCLUDED.%[3]s",
			sql_queries.UserIdColumnName,      // 1
			sql_queries.TokenValueColumnName,  // 2
			sql_queries.UpdatedAtColumnName)). // 3
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
