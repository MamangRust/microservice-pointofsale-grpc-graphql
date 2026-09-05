package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/microservice-point-of-sale-auth/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-point-of-sale-shared/errors"
)

type refreshTokenRepository struct {
	db *db.Queries
}

func NewRefreshTokenRepository(db *db.Queries) *refreshTokenRepository {
	return &refreshTokenRepository{
		db: db,
	}
}

func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*db.RefreshToken, error) {
	res, err := r.db.FindRefreshTokenByToken(ctx, token)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}

func (r *refreshTokenRepository) FindByUserId(ctx context.Context, user_id int) (*db.RefreshToken, error) {
	res, err := r.db.FindRefreshTokenByUserId(ctx, int32(user_id))
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}

func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*db.RefreshToken, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	res, err := r.db.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:     int32(req.UserId),
		Token:      req.Token,
		Expiration: expirationTime,
	})
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}

func (r *refreshTokenRepository) UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*db.RefreshToken, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	res, err := r.db.UpdateRefreshTokenByUserId(ctx, db.UpdateRefreshTokenByUserIdParams{
		UserID:     int32(req.UserId),
		Token:      req.Token,
		Expiration: expirationTime,
	})
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}

func (r *refreshTokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	err := r.db.DeleteRefreshToken(ctx, token)
	if err != nil {
		return sharedErrors.ErrInternal.WithInternal(err)
	}
	return nil
}

func (r *refreshTokenRepository) DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error {
	err := r.db.DeleteRefreshTokenByUserId(ctx, int32(user_id))
	if err != nil {
		return sharedErrors.ErrInternal.WithInternal(err)
	}
	return nil
}
