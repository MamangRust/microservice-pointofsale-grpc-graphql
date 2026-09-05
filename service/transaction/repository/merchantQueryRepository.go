package repository

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/adapter"
	"github.com/MamangRust/microservice-point-of-sale-pkg/resilience"
	db "github.com/MamangRust/microservice-point-of-sale-transacton/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/merchant_errors"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
)

type merchantQueryRepository struct {
	client pb.MerchantServiceClient
	guard  *resilience.DependencyGuard
}

func NewMerchantQueryRepository(client pb.MerchantServiceClient, opts ...adapter.GuardOption) MerchantQueryRepository {
	r := &merchantQueryRepository{
		client: client,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *merchantQueryRepository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func parseNullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// FindById resolves the merchant through the guarded gRPC dependency. The guard
// wraps ONLY the raw client call so transport-level failures trip the circuit
// breaker; business error conversion happens after the guard returns.
func (r *merchantQueryRepository) FindById(ctx context.Context, merchantID int) (*db.Merchant, error) {
	var resp *pb.ApiResponseMerchant
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindById(ctx, &pb.FindByIdMerchantRequest{
			Id: int32(merchantID),
		})
		return callErr
	})
	if err != nil {
		return nil, merchant_errors.ErrFindById
	}

	if resp == nil || resp.Data == nil {
		return nil, merchant_errors.ErrFindById
	}

	m := resp.Data
	return &db.Merchant{
		MerchantID:   m.Id,
		UserID:       m.UserId,
		Name:         m.Name,
		Description:  parseNullableString(m.Description),
		Address:      parseNullableString(m.Address),
		ContactEmail: parseNullableString(m.ContactEmail),
		ContactPhone: parseNullableString(m.ContactPhone),
		Status:       m.Status,
		CreatedAt:    parsePgTimestamp(m.CreatedAt),
		UpdatedAt:    parsePgTimestamp(m.UpdatedAt),
	}, nil
}
