package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-order/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/microservice-point-of-sale-shared/errors/order_item_errors"
)

type orderItemCommandRepository struct {
	db *db.Queries
}

func NewOrderItemCommandRepository(db *db.Queries) OrderItemCommandRepository {
	return &orderItemCommandRepository{
		db: db,
	}
}

func (r *orderItemCommandRepository) CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*db.OrderItem, error) {
	res, err := r.db.CreateOrderItem(ctx, db.CreateOrderItemParams{
		OrderID:   int32(req.OrderID),
		ProductID: int32(req.ProductID),
		Quantity:  int32(req.Quantity),
		Price:     int32(req.Price),
	})
	if err != nil {
		return nil, orderitem_errors.ErrCreateOrderItem
	}

	return res, nil
}

func (r *orderItemCommandRepository) UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*db.OrderItem, error) {
	res, err := r.db.UpdateOrderItem(ctx, db.UpdateOrderItemParams{
		OrderItemID: int32(req.OrderItemID),
		ProductID:   int32(req.ProductID),
		Quantity:    int32(req.Quantity),
		Price:       int32(req.Price),
	})
	if err != nil {
		return nil, orderitem_errors.ErrUpdateOrderItem
	}

	return res, nil
}

func (r *orderItemCommandRepository) DeleteOrderItem(ctx context.Context, order_id int) error {
	if err := r.db.DeleteOrderItem(ctx, int32(order_id)); err != nil {
		return orderitem_errors.ErrDeleteOrderItemPermanent
	}
	return nil
}
