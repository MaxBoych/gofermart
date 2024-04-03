package orderrepository

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/MaxBoych/gofermart/internal/order/ordermodels"
	"github.com/MaxBoych/gofermart/pkg/logger"
	"github.com/MaxBoych/gofermart/pkg/sqlqueries"
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

func (r *OrderRepo) GetOrder(ctx context.Context, number string) (*ordermodels.OrderStorageData, error) {
	query, args, err := sq.Select(sqlqueries.SelectOrder...).
		From(sqlqueries.OrderTableName).
		Where(sq.Eq{sqlqueries.OrderNumberColumnName: number}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Log.Error("Cannot to build sql SELECT query", zap.Error(err))
		return nil, err
	}

	data := ordermodels.OrderStorageData{}
	tr := r.txGetter.DefaultTrOrDB(ctx, r.db)
	err = tr.QueryRow(ctx, query, args...).Scan(&data.OrderID, &data.Number, &data.UserID, &data.Status, &data.Accrual, &data.CreatedAt, &data.UpdatedAt)
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

func (r *OrderRepo) CreateOrder(ctx context.Context, data ordermodels.OrderStorageData) error {
	query, args, err := sq.Insert(sqlqueries.OrderTableName).
		Columns(sqlqueries.InsertOrder...).
		Values(
			data.Number,
			data.UserID,
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

func (r *OrderRepo) GetOrders(ctx context.Context, userID int64) ([]ordermodels.OrderStorageData, error) {
	//var conditions sq.And
	//conditions = append(conditions, sq.Eq{sql_queries.UserIDColumnName: userID})
	//if nonFinalOnly {
	//	conditions = append(conditions, sq.NotEq{sql_queries.OrderStatusColumnName: order_models.OrderStatusInvalid})
	//	conditions = append(conditions, sq.NotEq{sql_queries.OrderStatusColumnName: order_models.OrderStatusProcessed})
	//}

	query, args, err := sq.Select(sqlqueries.SelectOrder...).
		From(sqlqueries.OrderTableName).
		Where(sq.Eq{sqlqueries.UserIDColumnName: userID}).
		OrderBy(sqlqueries.CreatedAtColumnName + " DESC").
		PlaceholderFormat(sq.Dollar).
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

	orders := make([]ordermodels.OrderStorageData, 0)
	for rows.Next() {
		order := ordermodels.OrderStorageData{}
		if err := rows.Scan(
			&order.OrderID,
			&order.Number,
			&order.UserID,
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

func (r *OrderRepo) UpdateOrder(ctx context.Context, updatedOrder ordermodels.OrderStorageData) error {
	query, args, err := sq.Update(sqlqueries.OrderTableName).
		Set(sqlqueries.OrderStatusColumnName, updatedOrder.Status).
		Set(sqlqueries.OrderAccrualColumnName, updatedOrder.Accrual).
		Where(sq.Eq{sqlqueries.OrderNumberColumnName: updatedOrder.Number}).
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
