package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	db "github.com/MamangRust/microservice-point-of-sale-order/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/order_errors"
)

type orderCommandRepository struct {
	db *db.Queries
}

func NewOrderCommandRepository(db *db.Queries) OrderCommandRepository {
	return &orderCommandRepository{
		db: db,
	}
}

func (r *orderCommandRepository) DeleteOrder(ctx context.Context, orderID int) error {
	if err := r.db.DeleteOrder(ctx, int32(orderID)); err != nil {
		return order_errors.ErrDeleteOrderPermanent
	}
	return nil
}

func (r *orderCommandRepository) CreateOrder(ctx context.Context, request *requests.CreateOrderRecordRequest) (*db.Order, error) {
	req := db.CreateOrderParams{
		MerchantID: int32(request.MerchantID),
		CashierID:  int32(request.CashierID),
		TotalPrice: int64(request.TotalPrice),
	}

	user, err := r.db.CreateOrder(ctx, req)
	if err != nil {
		return nil, order_errors.ErrCreateOrder
	}

	return user, nil
}

func (r *orderCommandRepository) FindAllTrashed(ctx context.Context) ([]*db.Order, error) {
	res, err := r.db.GetAllOrdersTrashed(ctx)
	if err != nil {
		return nil, order_errors.ErrFindAllTrashed
	}
	return res, nil
}

func (r *orderCommandRepository) UpdateOrder(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*db.Order, error) {
	req := db.UpdateOrderParams{
		OrderID:    int32(request.OrderID),
		TotalPrice: int64(request.TotalPrice),
	}

	res, err := r.db.UpdateOrder(ctx, req)
	if err != nil {
		return nil, order_errors.ErrUpdateOrder
	}

	return res, nil
}

func (r *orderCommandRepository) TrashedOrder(ctx context.Context, orderID int) (*db.Order, error) {
	res, err := r.db.TrashedOrder(ctx, int32(orderID))
	if err != nil {
		return nil, order_errors.ErrTrashedOrder.WithInternal(err)
	}

	return res, nil
}

func (r *orderCommandRepository) RestoreOrder(ctx context.Context, orderID int) (*db.Order, error) {
	res, err := r.db.RestoreOrder(ctx, int32(orderID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, order_errors.ErrRestoreOrderNotFound
		}
		return nil, order_errors.ErrRestoreOrder.WithInternal(err)
	}

	return res, nil
}

func (r *orderCommandRepository) DeleteOrderPermanent(ctx context.Context, orderID int) (bool, error) {
	_, err := r.db.DeleteOrderPermanently(ctx, int32(orderID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, order_errors.ErrDeleteOrderPermanentNotFound
		}
		return false, order_errors.ErrDeleteOrderPermanent.WithInternal(err)
	}

	return true, nil
}

func (r *orderCommandRepository) DeleteAllOrderPermanent(ctx context.Context) (bool, error) {
	_, err := r.db.DeleteAllPermanentOrders(ctx)
	if err != nil {
		return false, order_errors.ErrDeleteAllOrderPermanent.WithInternal(err)
	}
	return true, nil
}
