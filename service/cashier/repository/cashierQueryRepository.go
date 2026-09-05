package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-cashier/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/cashier_errors"
)

type cashierQueryRepository struct {
	db *db.Queries
}

func NewCashierQueryRepository(db *db.Queries) CashierQueryRepository {
	return &cashierQueryRepository{
		db: db,
	}
}

func (r *cashierQueryRepository) FindAllCashiers(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCashiersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCashiers(ctx, reqDb)
	if err != nil {
		return nil, nil, cashier_errors.ErrFindAllCashiers
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *cashierQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCashiersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCashiersActive(ctx, reqDb)
	if err != nil {
		return nil, nil, cashier_errors.ErrFindActiveCashiers
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *cashierQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCashiersTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCashiersTrashed(ctx, reqDb)
	if err != nil {
		return nil, nil, cashier_errors.ErrFindTrashedCashiers
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *cashierQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*db.GetCashiersByMerchantRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCashiersByMerchantParams{
		MerchantID: int32(req.MerchantID),
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetCashiersByMerchant(ctx, reqDb)
	if err != nil {
		return nil, nil, cashier_errors.ErrFindCashiersByMerchant
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *cashierQueryRepository) FindById(ctx context.Context, cashier_id int) (*db.Cashier, error) {
	res, err := r.db.GetCashierById(ctx, int32(cashier_id))
	if err != nil {
		return nil, cashier_errors.ErrFindCashierById
	}

	return res, nil
}
