package repository

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/adapter"
	"github.com/MamangRust/microservice-point-of-sale-pkg/resilience"
	db "github.com/MamangRust/microservice-point-of-sale-transacton/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/order_errors"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/jackc/pgx/v5/pgtype"
)

type orderQueryRepository struct {
	client pb.OrderServiceClient
	guard  *resilience.DependencyGuard
}

func NewOrderQueryRepository(client pb.OrderServiceClient, opts ...adapter.GuardOption) OrderQueryRepository {
	r := &orderQueryRepository{
		client: client,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *orderQueryRepository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

// FindById resolves the order through the guarded gRPC dependency. The guard
// wraps ONLY the raw client call so transport-level failures trip the circuit
// breaker; business error conversion happens after the guard returns.
func (r *orderQueryRepository) FindById(ctx context.Context, order_id int) (*db.Order, error) {
	var resp *pb.ApiResponseOrder
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindById(ctx, &pb.FindByIdOrderRequest{
			Id: int32(order_id),
		})
		return callErr
	})
	if err != nil {
		return nil, order_errors.ErrFindById
	}

	if resp == nil || resp.Data == nil {
		return nil, order_errors.ErrFindById
	}

	o := resp.Data
	return &db.Order{
		OrderID:    o.Id,
		MerchantID: o.MerchantId,
		CashierID:  o.CashierId,
		TotalPrice: int64(o.TotalPrice),
		CreatedAt:  parsePgTimestamp(o.CreatedAt),
		UpdatedAt:  parsePgTimestamp(o.UpdatedAt),
		DeletedAt:  pgtype.Timestamp{},
	}, nil
}
