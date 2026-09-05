package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-product/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/merchant_errors"
)

type merchantQueryRepository struct {
	db *db.Queries
}

func NewMerchantQueryRepository(db *db.Queries) *merchantQueryRepository {
	return &merchantQueryRepository{
		db: db,
	}
}

func (r *merchantQueryRepository) FindById(ctx context.Context, user_id int) (*db.Merchant, error) {
	res, err := r.db.GetMerchantByID(ctx, int32(user_id))
	if err != nil {
		return nil, merchant_errors.ErrFindById
	}

	return res, nil
}
