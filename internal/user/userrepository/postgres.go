package userrepository

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/MaxBoych/gofermart/internal/user/usermodels"
	"github.com/MaxBoych/gofermart/pkg/logger"
	"github.com/MaxBoych/gofermart/pkg/sqlqueries"
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

func (r *UserRepo) GetUserByLogin(ctx context.Context, login string) (*usermodels.UserStorageData, error) {
	query, args, err := sq.Select(sqlqueries.SelectUser...).
		From(sqlqueries.UserTableName).
		Where(sq.Eq{sqlqueries.LoginColumnName: login}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Log.Error("Cannot to build sql SELECT query", zap.Error(err))
		return nil, err
	}

	data := usermodels.UserStorageData{}
	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	err = tr.QueryRow(ctx, query, args...).Scan(&data.UserID, &data.Login, &data.Password, &data.CreatedAt, &data.UpdatedAt, &data.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			//logger.Log.Error("There is no such row", zap.Error(err))
			return nil, err
		}

		logger.Log.Error("Error while scanning, sql SELECT query", zap.Error(err))
		return nil, err
	}

	return &data, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, data usermodels.UserStorageData) (int64, error) {
	query, args, err := sq.Insert(sqlqueries.UserTableName).
		Columns(sqlqueries.InsertUser...).
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
	var userID int64
	err = tr.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		logger.Log.Error("Error while executing sql INSERT query", zap.Error(err))
		return -1, err
	}

	return userID, nil
}
