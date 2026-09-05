package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-order/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/order_errors"
)

type orderQueryRepository struct {
	db *db.Queries
}

func NewOrderQueryRepository(db *db.Queries) OrderQueryRepository {
	return &orderQueryRepository{
		db: db,
	}
}

func (r *orderQueryRepository) FindAllOrders(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrders(ctx, reqDb)
	if err != nil {
		return nil, nil, order_errors.ErrFindAllOrders
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	return res, &totalCount, nil
}

func (r *orderQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersActive(ctx, reqDb)
	if err != nil {
		return nil, nil, order_errors.ErrFindByActive
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	return res, &totalCount, nil
}

func (r *orderQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersTrashed(ctx, reqDb)
	if err != nil {
		return nil, nil, order_errors.ErrFindByTrashed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	return res, &totalCount, nil
}

func (r *orderQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllOrderMerchant) ([]*db.GetOrdersByMerchantRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersByMerchantParams{
		Column1: req.Search,
		Column4: int32(req.MerchantID),
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersByMerchant(ctx, reqDb)
	if err != nil {
		return nil, nil, order_errors.ErrFindByMerchant
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	return res, &totalCount, nil
}

func (r *orderQueryRepository) FindById(ctx context.Context, orderID int) (*db.Order, error) {
	res, err := r.db.GetOrderByID(ctx, int32(orderID))
	if err != nil {
		return nil, order_errors.ErrFindById
	}

	return res, nil
}

func (r *orderQueryRepository) FindByTrashedId(ctx context.Context, orderID int) (*db.Order, error) {
	res, err := r.db.GetOrderByIDTrashed(ctx, int32(orderID))
	if err != nil {
		return nil, order_errors.ErrFindByTrashedId
	}

	return res, nil
}
