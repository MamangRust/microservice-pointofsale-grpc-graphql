package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-product/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/product_errors"
)

type productQueryRepository struct {
	db *db.Queries
}

func NewProductQueryRepository(db *db.Queries) *productQueryRepository {
	return &productQueryRepository{
		db: db,
	}
}

func (r *productQueryRepository) FindAllProducts(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProducts(ctx, reqDb)
	if err != nil {
		return nil, nil, product_errors.ErrFindAllProducts
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *productQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsActive(ctx, reqDb)
	if err != nil {
		return nil, nil, product_errors.ErrFindByActive
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *productQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsTrashed(ctx, reqDb)
	if err != nil {
		return nil, nil, product_errors.ErrFindByTrashed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *productQueryRepository) FindByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*db.GetProductsByMerchantRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsByMerchantParams{
		MerchantID: int32(req.MerchantID),
		Column2:    &req.Search,
		Column3:    req.CategoryID,
		Column4:    req.MinPrice,
		Column5:    req.MaxPrice,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetProductsByMerchant(ctx, reqDb)
	if err != nil {
		return nil, nil, product_errors.ErrFindByMerchant
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *productQueryRepository) FindByCategory(ctx context.Context, req *requests.ProductByCategoryRequest) ([]*db.GetProductsByCategoryNameRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsByCategoryNameParams{
		Name:    req.CategoryName,
		Column2: req.Search,
		Column3: int32MinPrice(req.MinPrice),
		Column4: int32MaxPrice(req.MaxPrice),
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsByCategoryName(ctx, reqDb)
	if err != nil {
		return nil, nil, product_errors.ErrFindByCategory.WithInternal(err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

// int32MinPrice maps an optional minimum-price filter to the int32 column
// expected by the generated query. nil (and 0) mean "no lower bound".
func int32MinPrice(p *int) int32 {
	if p == nil {
		return 0
	}
	return int32(*p)
}

// int32MaxPrice maps an optional maximum-price filter to the int32 column
// expected by the generated query. nil (and 0) mean "no upper bound" and map
// to a very high sentinel, mirroring the GetProductsByMerchant query.
func int32MaxPrice(p *int) int32 {
	if p == nil || *p == 0 {
		return 999999999
	}
	return int32(*p)
}

func (r *productQueryRepository) FindById(ctx context.Context, product_id int) (*db.Product, error) {
	res, err := r.db.GetProductByID(ctx, int32(product_id))
	if err != nil {
		return nil, product_errors.ErrFindById
	}

	return res, nil
}

func (r *productQueryRepository) FindByIdTrashed(ctx context.Context, product_id int) (*db.Product, error) {
	res, err := r.db.GetProductByIdTrashed(ctx, int32(product_id))
	if err != nil {
		return nil, product_errors.ErrFindByIdTrashed
	}

	return res, nil
}
