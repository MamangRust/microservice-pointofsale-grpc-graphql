package repository

import (
	"context"
	"errors"

	db "github.com/MamangRust/microservice-point-of-sale-order/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/order_errors"
	"github.com/jackc/pgx/v5"
)

type productCommandRepository struct {
	db *db.Queries
}

func NewProductCommandRepository(db *db.Queries) ProductCommandRepository {
	return &productCommandRepository{
		db: db,
	}
}

func (r *productCommandRepository) DecrementProductCountStock(ctx context.Context, productID int, quantity int) (*db.Product, error) {
	res, err := r.db.DecrementProductCountStock(ctx, db.DecrementProductCountStockParams{
		ProductID: int32(productID),
		Quantity:  int32(quantity),
	})
	if err != nil {
		// Atomic guard `count_in_stock >= quantity` returned no rows:
		// stock is insufficient (never negative) → 400 Bad Request.
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, order_errors.ErrInsufficientProductStock
		}
		return nil, err
	}

	return res, nil
}

func (r *productCommandRepository) IncrementProductCountStock(ctx context.Context, productID int, quantity int) (*db.Product, error) {
	res, err := r.db.IncrementProductCountStock(ctx, db.IncrementProductCountStockParams{
		ProductID: int32(productID),
		Quantity:  int32(quantity),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
