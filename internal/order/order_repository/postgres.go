package order_repository

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/MaxBoych/gofermart/internal/order/order_models"
	"github.com/MaxBoych/gofermart/pkg/logger"
	"github.com/MaxBoych/gofermart/pkg/sql_queries"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type OrderRepo struct {
	db       *pgxpool.Pool
	txGetter *trmpgx.CtxGetter
}

func NewOrderRepo(db *pgxpool.Pool, txGetter *trmpgx.CtxGetter) *OrderRepo {
	return &OrderRepo{
		db:       db,
		txGetter: txGetter,
	}
}

func (r *OrderRepo) GetOrder(ctx context.Context, number string) (*order_models.OrderStorageData, error) {
	query, args, err := sq.Select(sql_queries.SelectOrder...).
		From(sql_queries.OrderTableName).
		Where(sq.Eq{sql_queries.OrderNumberColumnName: number}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Log.Error("Cannot to build sql SELECT query", zap.Error(err))
		return nil, err
	}

	data := order_models.OrderStorageData{}
	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	err = tr.QueryRow(ctx, query, args...).Scan(&data.OrderId, &data.Number, &data.UserId, &data.Status, &data.Accrual, &data.CreatedAt, &data.UpdatedAt)
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

func (r *OrderRepo) CreateOrder(ctx context.Context, data order_models.OrderStorageData) error {
	query, args, err := sq.Insert(sql_queries.OrderTableName).
		Columns(sql_queries.InsertOrder...).
		Values(
			data.Number,
			data.UserId,
			data.Status,
			data.Accrual,
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

func (r *OrderRepo) GetOrders(ctx context.Context, userId int64) ([]order_models.OrderStorageData, error) {
	//var conditions sq.And
	//conditions = append(conditions, sq.Eq{sql_queries.UserIdColumnName: userId})
	//if nonFinalOnly {
	//	conditions = append(conditions, sq.NotEq{sql_queries.OrderStatusColumnName: order_models.OrderStatusInvalid})
	//	conditions = append(conditions, sq.NotEq{sql_queries.OrderStatusColumnName: order_models.OrderStatusProcessed})
	//}

	query, args, err := sq.Select(sql_queries.SelectOrder...).
		From(sql_queries.OrderTableName).
		Where(sq.Eq{sql_queries.UserIdColumnName: userId}).
		OrderBy(sql_queries.CreatedAtColumnName + " DESC").
		ToSql()
	if err != nil {
		logger.Log.Error("Error to make sql builder, sql SELECT query", zap.Error(err))
		return nil, err
	}

	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	rows, err := tr.Query(ctx, query, args...)
	if err != nil {
		logger.Log.Error("Error to execute sql SELECT query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	orders := make([]order_models.OrderStorageData, 0)
	for rows.Next() {
		order := order_models.OrderStorageData{}
		if err := rows.Scan(
			&order.OrderId,
			&order.Number,
			&order.UserId,
			&order.Status,
			&order.Accrual,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			logger.Log.Error("Error to scan rows", zap.Error(err))
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepo) UpdateOrder(ctx context.Context, updatedOrder order_models.OrderStorageData) error {
	query, args, err := sq.Update(sql_queries.OrderTableName).
		Set(sql_queries.OrderStatusColumnName, updatedOrder.Status).
		Set(sql_queries.OrderAccrualColumnName, updatedOrder.Accrual).
		Where(sq.Eq{sql_queries.OrderNumberColumnName: updatedOrder.Number}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		logger.Log.Error("Error to build sql UPDATE query", zap.Error(err))
		return err
	}

	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	_, err = tr.Exec(ctx, query, args...)
	if err != nil {
		logger.Log.Error("Error to execute sql UPDATE query", zap.Error(err))
		return err
	}

	return nil
}
