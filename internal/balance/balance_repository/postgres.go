package balance_repository

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/MaxBoych/gofermart/internal/balance/balance_models"
	"github.com/MaxBoych/gofermart/pkg/logger"
	"github.com/MaxBoych/gofermart/pkg/sql_queries"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type BalanceRepo struct {
	db       *pgxpool.Pool
	txGetter *trmpgx.CtxGetter
}

func NewBalanceRepo(db *pgxpool.Pool, txGetter *trmpgx.CtxGetter) *BalanceRepo {
	return &BalanceRepo{
		db:       db,
		txGetter: txGetter,
	}
}

func (r *BalanceRepo) GetBalance(ctx context.Context, userId int64) (*balance_models.BalanceStorageData, error) {
	query, args, err := sq.Select(sql_queries.SelectBalance...).
		From(sql_queries.BalanceTableName).
		Where(sq.Eq{sql_queries.UserIdColumnName: userId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Log.Error("Cannot to build sql SELECT query", zap.Error(err))
		return nil, err
	}

	data := balance_models.BalanceStorageData{}
	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	err = tr.QueryRow(ctx, query, args...).Scan(
		&data.BalanceId,
		&data.UserId,
		&data.Current,
		&data.Withdrawn,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
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

func (r *BalanceRepo) CreateBalance(ctx context.Context, userId int64) error {
	query, args, err := sq.Insert(sql_queries.BalanceTableName).
		Columns(sql_queries.InsertBalance...).
		Values(
			userId,
			0,
			0,
			time.Now(),
			time.Now(),
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Log.Error("Cannot to build sql INSERT query", zap.Error(err))
		return err
	}

	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	_, err = tr.Exec(ctx, query, args...)
	if err != nil {
		logger.Log.Error("Error while executing sql INSERT query", zap.Error(err))
		return err
	}

	return nil
}

func (r *BalanceRepo) UpdateBalance(ctx context.Context, req balance_models.BalanceChangeData) error {
	updateBuilder := sq.Update(sql_queries.BalanceTableName).
		Set(sql_queries.BalanceCurrentColumnName, sq.Expr(sql_queries.BalanceCurrentColumnName+" "+req.Action+" ?", req.Sum))

	if req.IsWithdraw() {
		updateBuilder = updateBuilder.Set(sql_queries.BalanceWithdrawnColumnName, sq.Expr(sql_queries.BalanceWithdrawnColumnName+" + 1"))
	}

	query, args, err := updateBuilder.
		Where(sq.Eq{sql_queries.UserIdColumnName: req.UserId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Log.Error("Error to build sql UPDATE query", zap.Error(err))
		return err
	}
	logger.Log.Info("UpdateBalance SQL query", zap.String("query", query))

	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	_, err = tr.Exec(ctx, query, args...)
	if err != nil {
		logger.Log.Error("Error to execute sql UPDATE query", zap.Error(err))
		return err
	}

	return nil
}

func (r *BalanceRepo) CreateWithdraw(ctx context.Context, req balance_models.WithdrawRequestData) error {
	query, args, err := sq.Insert(sql_queries.WithdrawTableName).
		Columns(sql_queries.InsertWithdraw...).
		Values(
			req.Order,
			req.Sum,
			req.UserId,
			time.Now(),
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Log.Error("Error to build sql INSERT query", zap.Error(err))
		return err
	}

	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	_, err = tr.Exec(ctx, query, args...)
	if err != nil {
		logger.Log.Error("Error while executing sql INSERT query", zap.Error(err))
		return err
	}

	return nil
}

func (r *BalanceRepo) GetWithdrawals(ctx context.Context, userId int64) ([]balance_models.WithdrawStorageData, error) {
	query, args, err := sq.Select(sql_queries.SelectWithdraw...).
		From(sql_queries.WithdrawTableName).
		Where(sq.Eq{sql_queries.UserIdColumnName: userId}).
		OrderBy(sql_queries.CreatedAtColumnName + " DESC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Log.Error("Error to make sql builder, sql SELECT query", zap.Error(err))
		return nil, err
	}
	logger.Log.Info("SelectWithdrawals SQL query", zap.String("query", query))

	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	rows, err := tr.Query(ctx, query, args...)
	if err != nil {
		logger.Log.Error("Error to execute sql SELECT query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	withdrawals := make([]balance_models.WithdrawStorageData, 0)
	for rows.Next() {
		withdraw := balance_models.WithdrawStorageData{}
		if err := rows.Scan(
			&withdraw.WithdrawId,
			&withdraw.Order,
			&withdraw.Sum,
			&withdraw.UserId,
			&withdraw.CreatedAt,
		); err != nil {
			logger.Log.Error("Error to scan rows", zap.Error(err))
			return nil, err
		}
		withdrawals = append(withdrawals, withdraw)
	}

	return withdrawals, nil
}
