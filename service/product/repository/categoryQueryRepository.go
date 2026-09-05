package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-product/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/category_errors"
)

type categoryQueryRepository struct {
	db *db.Queries
}

func NewCategoryQueryRepository(db *db.Queries) *categoryQueryRepository {
	return &categoryQueryRepository{
		db: db,
	}
}

func (r *categoryQueryRepository) FindById(ctx context.Context, category_id int) (*db.Category, error) {
	res, err := r.db.GetCategoryByID(ctx, int32(category_id))
	if err != nil {
		return nil, category_errors.ErrFindById
	}

	return res, nil
}
