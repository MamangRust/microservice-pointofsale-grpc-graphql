package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-product/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
)

type CategoryQueryRepository interface {
	FindById(ctx context.Context, category_id int) (*db.Category, error)
}

type MerchantQueryRepository interface {
	FindById(ctx context.Context, merchant_id int) (*db.Merchant, error)
}

type ProductQueryRepository interface {
	FindAllProducts(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsTrashedRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*db.GetProductsByMerchantRow, *int, error)
	FindByCategory(ctx context.Context, req *requests.ProductByCategoryRequest) ([]*db.GetProductsByCategoryNameRow, *int, error)
	FindById(ctx context.Context, product_id int) (*db.Product, error)
	FindByIdTrashed(ctx context.Context, id int) (*db.Product, error)
}

type ProductCommandRepository interface {
	CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*db.Product, error)
	UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*db.Product, error)
	UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*db.Product, error)
	TrashedProduct(ctx context.Context, product_id int) (*db.Product, error)
	RestoreProduct(ctx context.Context, product_id int) (*db.Product, error)
	DeleteProductPermanent(ctx context.Context, product_id int) (bool, error)
	RestoreAllProducts(ctx context.Context) (bool, error)
	DeleteAllProductPermanent(ctx context.Context) (bool, error)
}
