package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-category/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/category_errors"
)

type categoryQueryRepository struct {
	db *db.Queries
}

func NewCategoryQueryRepository(db *db.Queries) CategoryQueryRepository {
	return &categoryQueryRepository{
		db: db,
	}
}

func (r *categoryQueryRepository) FindAllCategory(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCategoriesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategories(ctx, reqDb)
	if err != nil {
		return nil, nil, category_errors.ErrFindAllCategory
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *categoryQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCategoriesActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategoriesActive(ctx, reqDb)
	if err != nil {
		return nil, nil, category_errors.ErrFindByActive
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *categoryQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCategoriesTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategoriesTrashed(ctx, reqDb)
	if err != nil {
		return nil, nil, category_errors.ErrFindByTrashed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *categoryQueryRepository) FindById(ctx context.Context, category_id int) (*db.Category, error) {
	res, err := r.db.GetCategoryByID(ctx, int32(category_id))
	if err != nil {
		return nil, category_errors.ErrFindById
	}

	return res, nil
}

func (r *categoryQueryRepository) FindByIdTrashed(ctx context.Context, category_id int) (*db.Category, error) {
	res, err := r.db.GetCategoryByIDTrashed(ctx, int32(category_id))
	if err != nil {
		return nil, category_errors.ErrFindById
	}

	return res, nil
}

func (r *categoryQueryRepository) FindByName(ctx context.Context, name string) (*db.Category, error) {
	res, err := r.db.GetCategoryByName(ctx, name)
	if err != nil {
		return nil, category_errors.ErrFindByName
	}

	return res, nil
}

func (r *categoryQueryRepository) FindByNameAndId(ctx context.Context, req *requests.CategoryNameAndId) (*db.Category, error) {
	res, err := r.db.GetCategoryByNameAndId(ctx, db.GetCategoryByNameAndIdParams{
		Name:       req.Name,
		CategoryID: int32(req.CategoryID),
	})
	if err != nil {
		return nil, category_errors.ErrFindByNameAndId
	}

	return res, nil
}
