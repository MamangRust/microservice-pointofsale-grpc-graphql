package repository

import (
	"context"
	"errors"

	db "github.com/MamangRust/microservice-point-of-sale-user/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/user_errors"
	"github.com/jackc/pgx/v5"
)

type userQueryRepository struct {
	db *db.Queries
}

func NewUserQueryRepository(db *db.Queries) *userQueryRepository {
	return &userQueryRepository{
		db: db,
	}
}

func (r *userQueryRepository) FindAllUsers(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetUsersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetUsers(ctx, reqDb)
	if err != nil {
		return nil, nil, user_errors.ErrFindAllUsers
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *userQueryRepository) FindById(ctx context.Context, user_id int) (*db.User, error) {
	res, err := r.db.GetUserByID(ctx, int32(user_id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrInternalServerError.WithInternal(err)
	}

	return res, nil
}

func (r *userQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetUsersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetUsersActive(ctx, reqDb)
	if err != nil {
		return nil, nil, user_errors.ErrFindActiveUsers
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *userQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUserTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetUserTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetUserTrashed(ctx, reqDb)
	if err != nil {
		return nil, nil, user_errors.ErrFindTrashedUsers
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *userQueryRepository) FindByEmail(ctx context.Context, email string) (*db.User, error) {
	res, err := r.db.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, user_errors.ErrInternalServerError.WithInternal(err)
	}

	return res, nil
}
