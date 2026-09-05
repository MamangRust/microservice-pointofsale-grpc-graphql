package repository

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/adapter"
	"github.com/MamangRust/microservice-point-of-sale-pkg/resilience"
	db "github.com/MamangRust/microservice-point-of-sale-order/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/product_errors"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
)

type productQueryRepository struct {
	client pb.ProductServiceClient
	guard  *resilience.DependencyGuard
}

func NewProductQueryRepository(client pb.ProductServiceClient, opts ...adapter.GuardOption) ProductQueryRepository {
	r := &productQueryRepository{
		client: client,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *productQueryRepository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func parseNullableInt32(i int32) *int32 {
	return &i
}

// FindById resolves the product through the guarded gRPC dependency. The guard
// wraps ONLY the raw client call so transport-level failures trip the circuit
// breaker; business error conversion happens after the guard returns.
func (r *productQueryRepository) FindById(ctx context.Context, product_id int) (*db.Product, error) {
	var resp *pb.ApiResponseProduct
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindById(ctx, &pb.FindByIdProductRequest{
			Id: int32(product_id),
		})
		return callErr
	})
	if err != nil {
		return nil, product_errors.ErrFindById
	}

	if resp == nil || resp.Data == nil {
		return nil, product_errors.ErrFindById
	}

	p := resp.Data
	return &db.Product{
		ProductID:    p.Id,
		MerchantID:   p.MerchantId,
		CategoryID:   p.CategoryId,
		Name:         p.Name,
		Description:  parseNullableString(p.Description),
		Price:        p.Price,
		CountInStock: p.CountInStock,
		Brand:        parseNullableString(p.Brand),
		Weight:       parseNullableInt32(p.Weight),
		SlugProduct:  parseNullableString(p.SlugProduct),
		ImageProduct: parseNullableString(p.ImageProduct),
		Barcode:      parseNullableString(p.Barcode),
		CreatedAt:    parsePgTimestamp(p.CreatedAt),
		UpdatedAt:    parsePgTimestamp(p.UpdatedAt),
	}, nil
}
