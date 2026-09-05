package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-order-item/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
)

type orderItemQueryRepository struct {
	db *db.Queries
}

func NewOrderItemQueryRepository(db *db.Queries) *orderItemQueryRepository {
	return &orderItemQueryRepository{
		db: db,
	}
}

func (r *orderItemQueryRepository) FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItems(ctx, reqDb)
	if err != nil {
		return nil, nil, err
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *orderItemQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItemsActive(ctx, reqDb)
	if err != nil {
		return nil, nil, err
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *orderItemQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItemsTrashed(ctx, reqDb)
	if err != nil {
		return nil, nil, err
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, order_id int) ([]*db.OrderItem, error) {
	res, err := r.db.GetOrderItemsByOrder(ctx, int32(order_id))
	if err != nil {
		return nil, err
	}

	return res, nil
}
