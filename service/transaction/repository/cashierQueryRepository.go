package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-point-of-sale-pkg/adapter"
	"github.com/MamangRust/microservice-point-of-sale-pkg/resilience"
	db "github.com/MamangRust/microservice-point-of-sale-transacton/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/cashier_errors"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/jackc/pgx/v5/pgtype"
)

type cashierQueryRepository struct {
	client pb.CashierServiceClient
	guard  *resilience.DependencyGuard
}

func NewCashierQueryRepository(client pb.CashierServiceClient, opts ...adapter.GuardOption) CashierQueryRepository {
	r := &cashierQueryRepository{
		client: client,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *cashierQueryRepository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func parsePgTimestamp(s string) pgtype.Timestamp {
	if s == "" {
		return pgtype.Timestamp{}
	}
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return pgtype.Timestamp{}
		}
	}
	return pgtype.Timestamp{Time: t, Valid: true}
}

// FindById resolves the cashier through the guarded gRPC dependency. The guard
// wraps ONLY the raw client call so transport-level failures trip the circuit
// breaker; business error conversion happens after the guard returns.
func (r *cashierQueryRepository) FindById(ctx context.Context, cashier_id int) (*db.Cashier, error) {
	var resp *pb.ApiResponseCashier
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindById(ctx, &pb.FindByIdCashierRequest{
			Id: int32(cashier_id),
		})
		return callErr
	})
	if err != nil {
		return nil, cashier_errors.ErrFindCashierById
	}

	if resp == nil || resp.Data == nil {
		return nil, cashier_errors.ErrFindCashierById
	}

	c := resp.Data
	return &db.Cashier{
		CashierID:  c.Id,
		MerchantID: c.MerchantId,
		Name:       c.Name,
		CreatedAt:  parsePgTimestamp(c.CreatedAt),
		UpdatedAt:  parsePgTimestamp(c.UpdatedAt),
	}, nil
}
