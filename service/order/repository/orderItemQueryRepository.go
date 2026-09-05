package repository

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/adapter"
	"github.com/MamangRust/microservice-point-of-sale-pkg/resilience"
	db "github.com/MamangRust/microservice-point-of-sale-order/database/schema"
	orderitem_errors "github.com/MamangRust/microservice-point-of-sale-shared/errors/order_item_errors"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
)

type orderItemQueryRepository struct {
	client pb.OrderItemServiceClient
	guard  *resilience.DependencyGuard
}

func NewOrderItemQueryRepository(client pb.OrderItemServiceClient, opts ...adapter.GuardOption) OrderItemQueryRepository {
	r := &orderItemQueryRepository{
		client: client,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *orderItemQueryRepository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *orderItemQueryRepository) CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error) {
	items, err := r.FindOrderItemByOrder(ctx, orderID)
	if err != nil {
		return nil, orderitem_errors.ErrCalculateTotalPrice
	}

	var total int32 = 0
	for _, item := range items {
		if item != nil {
			total += item.Quantity * item.Price
		}
	}

	return &total, nil
}

// FindOrderItemByOrder resolves order items through the guarded gRPC
// dependency. The guard wraps ONLY the raw client call so transport-level
// failures trip the circuit breaker; business error conversion happens after
// the guard returns.
func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, orderID int) ([]*db.OrderItem, error) {
	var resp *pb.ApiResponsesOrderItem
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindOrderItemByOrder(ctx, &pb.FindByIdOrderItemRequest{
			Id: int32(orderID),
		})
		return callErr
	})
	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}

	if resp == nil || resp.Data == nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}

	var res []*db.OrderItem
	for _, item := range resp.Data {
		if item == nil {
			continue
		}
		res = append(res, &db.OrderItem{
			OrderItemID: item.Id,
			OrderID:     item.OrderId,
			ProductID:   item.ProductId,
			Quantity:    item.Quantity,
			Price:       item.Price,
			CreatedAt:    parsePgTimestamp(item.CreatedAt),
			UpdatedAt:    parsePgTimestamp(item.UpdatedAt),
		})
	}

	return res, nil
}
