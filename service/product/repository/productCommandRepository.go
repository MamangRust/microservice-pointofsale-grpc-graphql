package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-product/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/product_errors"
)

type productCommandRepository struct {
	db *db.Queries
}

func NewProductCommandRepository(db *db.Queries) *productCommandRepository {
	return &productCommandRepository{
		db: db,
	}
}

func (r *productCommandRepository) CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*db.Product, error) {
	weightVal := int32(request.Weight)
	req := db.CreateProductParams{
		MerchantID:   int32(request.MerchantID),
		CategoryID:   int32(request.CategoryID),
		Name:         request.Name,
		Description:  &request.Description,
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        &request.Brand,
		Weight:       &weightVal,
		SlugProduct:  request.SlugProduct,
		ImageProduct: &request.ImageProduct,
	}

	product, err := r.db.CreateProduct(ctx, req)
	if err != nil {
		return nil, product_errors.ErrCreateProduct
	}

	return product, nil
}

func (r *productCommandRepository) UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*db.Product, error) {
	weightVal := int32(request.Weight)
	req := db.UpdateProductParams{
		ProductID:    int32(*request.ProductID),
		CategoryID:   int32(request.CategoryID),
		Name:         request.Name,
		Description:  &request.Description,
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        &request.Brand,
		Weight:       &weightVal,
		ImageProduct: &request.ImageProduct,
	}

	res, err := r.db.UpdateProduct(ctx, req)
	if err != nil {
		return nil, product_errors.ErrUpdateProduct
	}

	return res, nil
}

func (r *productCommandRepository) UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*db.Product, error) {
	res, err := r.db.UpdateProductCountStock(ctx, db.UpdateProductCountStockParams{
		ProductID:    int32(product_id),
		CountInStock: int32(stock),
	})
	if err != nil {
		return nil, product_errors.ErrUpdateProductCountStock
	}

	return res, nil
}

func (r *productCommandRepository) TrashedProduct(ctx context.Context, product_id int) (*db.Product, error) {
	res, err := r.db.TrashProduct(ctx, int32(product_id))
	if err != nil {
		return nil, product_errors.ErrTrashedProduct
	}

	return res, nil
}

func (r *productCommandRepository) RestoreProduct(ctx context.Context, product_id int) (*db.Product, error) {
	res, err := r.db.RestoreProduct(ctx, int32(product_id))
	if err != nil {
		return nil, product_errors.ErrRestoreProduct
	}

	return res, nil
}

func (r *productCommandRepository) DeleteProductPermanent(ctx context.Context, product_id int) (bool, error) {
	err := r.db.DeleteProductPermanently(ctx, int32(product_id))
	if err != nil {
		return false, product_errors.ErrDeleteProductPermanent
	}

	return true, nil
}

func (r *productCommandRepository) RestoreAllProducts(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllProducts(ctx)
	if err != nil {
		return false, product_errors.ErrRestoreAllProducts
	}

	return true, nil
}

func (r *productCommandRepository) DeleteAllProductPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentProducts(ctx)
	if err != nil {
		return false, product_errors.ErrDeleteAllProductPermanent
	}

	return true, nil
}
