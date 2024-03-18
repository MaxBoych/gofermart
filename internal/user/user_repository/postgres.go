package user_repository

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/MaxBoych/gofermart/internal/store/sql_queries"
	"github.com/MaxBoych/gofermart/internal/user/user_models"
	"github.com/MaxBoych/gofermart/pkg/logger"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type UserRepo struct {
	db       *pgxpool.Pool
	txGetter *trmpgx.CtxGetter
}

func NewUserRepo(db *pgxpool.Pool, txGetter *trmpgx.CtxGetter) *UserRepo {
	return &UserRepo{
		db:       db,
		txGetter: txGetter,
	}
}

func (r *UserRepo) GetUserByLogin(ctx context.Context, login string) (*user_models.UserStorageData, error) {
	query, args, err := sq.Select(sql_queries.SelectUser...).
		From(sql_queries.UserTableName).
		Where(sq.Eq{sql_queries.LoginColumnName: login}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Log.Error("Cannot to build sql SELECT query", zap.Error(err))
		return nil, err
	}

	data := user_models.UserStorageData{}
	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	err = tr.QueryRow(ctx, query, args...).Scan(&data.UserId, &data.Login, &data.Password, &data.CreatedAt, &data.UpdatedAt, &data.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Log.Error("There is no such row", zap.Error(err))
			return nil, err
		}

		logger.Log.Error("Error while scanning, sql SELECT query", zap.Error(err))
		return nil, err
	}

	return &data, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, data user_models.UserStorageData) (int64, error) {
	query, args, err := sq.Insert(sql_queries.UserTableName).
		Columns(sql_queries.InsertUser...).
		Values(
			data.Login,
			data.Password,
			time.Now(),
			time.Now(),
		).
		Suffix("RETURNING user_id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Log.Error("Cannot to build sql INSERT query", zap.Error(err))
		return -1, err
	}

	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	var userId int64
	err = tr.QueryRow(ctx, query, args...).Scan(&userId)
	if err != nil {
		logger.Log.Error("Error while executing sql INSERT query", zap.Error(err))
		return -1, err
	}

	return userId, nil
}
