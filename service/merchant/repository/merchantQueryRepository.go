package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-merchant/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-point-of-sale-shared/errors"
)

type merchantQueryRepository struct {
	db *db.Queries
}

func NewMerchantQueryRepository(db *db.Queries) MerchantQueryRepository {
	return &merchantQueryRepository{
		db: db,
	}
}

func (r *merchantQueryRepository) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchants(ctx, reqDb)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *merchantQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsActive(ctx, reqDb)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *merchantQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsTrashed(ctx, reqDb)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *merchantQueryRepository) FindById(ctx context.Context, userID int) (*db.Merchant, error) {
	res, err := r.db.GetMerchantByID(ctx, int32(userID))
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}
