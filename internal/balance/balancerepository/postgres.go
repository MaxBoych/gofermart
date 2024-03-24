package balancerepository

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/MaxBoych/gofermart/internal/balance/balancemodels"
	"github.com/MaxBoych/gofermart/pkg/logger"
	"github.com/MaxBoych/gofermart/pkg/sqlqueries"
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

func (r *BalanceRepo) GetBalance(ctx context.Context, userID int64) (*balancemodels.BalanceStorageData, error) {
	query, args, err := sq.Select(sqlqueries.SelectBalance...).
		From(sqlqueries.BalanceTableName).
		Where(sq.Eq{sqlqueries.UserIDColumnName: userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Log.Error("Cannot to build sql SELECT query", zap.Error(err))
		return nil, err
	}

	data := balancemodels.BalanceStorageData{}
	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	err = tr.QueryRow(ctx, query, args...).Scan(
		&data.BalanceID,
		&data.UserID,
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

func (r *BalanceRepo) CreateBalance(ctx context.Context, userID int64) error {
	query, args, err := sq.Insert(sqlqueries.BalanceTableName).
		Columns(sqlqueries.InsertBalance...).
		Values(
			userID,
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

func (r *BalanceRepo) UpdateBalance(ctx context.Context, changeData balancemodels.BalanceChangeData) error {
	updateBuilder := sq.Update(sqlqueries.BalanceTableName).
		Set(sqlqueries.BalanceCurrentColumnName, sq.Expr(sqlqueries.BalanceCurrentColumnName+" "+changeData.Action+" ?", changeData.Sum))

	if changeData.IsWithdraw() {
		updateBuilder = updateBuilder.Set(sqlqueries.BalanceWithdrawnColumnName, sq.Expr(sqlqueries.BalanceWithdrawnColumnName+" + ?", changeData.Sum))
	}

	query, args, err := updateBuilder.
		Where(sq.Eq{sqlqueries.UserIDColumnName: changeData.UserID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Log.Error("Error to build sql UPDATE query", zap.Error(err))
		return err
	}
	logger.Log.Info("UpdateBalance SQL query", zap.String("query", query))
	for _, a := range args {
		logger.Log.Info("UpdateBalance SQL args", zap.Any("arg", a))
	}

	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	_, err = tr.Exec(ctx, query, args...)
	if err != nil {
		logger.Log.Error("Error to execute sql UPDATE query", zap.Error(err))
		return err
	}

	return nil
}

func (r *BalanceRepo) CreateWithdraw(ctx context.Context, req balancemodels.WithdrawRequestData) error {
	query, args, err := sq.Insert(sqlqueries.WithdrawTableName).
		Columns(sqlqueries.InsertWithdraw...).
		Values(
			req.Order,
			req.Sum,
			req.UserID,
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

func (r *BalanceRepo) GetWithdrawals(ctx context.Context, userId int64) ([]balancemodels.WithdrawStorageData, error) {
	query, args, err := sq.Select(sqlqueries.SelectWithdraw...).
		From(sqlqueries.WithdrawTableName).
		Where(sq.Eq{sqlqueries.UserIDColumnName: userId}).
		OrderBy(sqlqueries.CreatedAtColumnName + " DESC").
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

	withdrawals := make([]balancemodels.WithdrawStorageData, 0)
	for rows.Next() {
		withdraw := balancemodels.WithdrawStorageData{}
		if err := rows.Scan(
			&withdraw.WithdrawID,
			&withdraw.Order,
			&withdraw.Sum,
			&withdraw.UserID,
			&withdraw.CreatedAt,
		); err != nil {
			logger.Log.Error("Error to scan rows", zap.Error(err))
			return nil, err
		}
		withdrawals = append(withdrawals, withdraw)
	}

	return withdrawals, nil
}
