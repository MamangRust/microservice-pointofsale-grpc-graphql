package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/microservice-point-of-sale-auth/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-point-of-sale-shared/errors"
	"github.com/jackc/pgx/v5"
)

type resetTokenRepository struct {
	db *db.Queries
}

func NewResetTokenRepository(db *db.Queries) *resetTokenRepository {
	return &resetTokenRepository{
		db: db,
	}
}

func (r *resetTokenRepository) FindByToken(ctx context.Context, code string) (*db.ResetToken, error) {
	res, err := r.db.GetResetToken(ctx, code)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}

func (r *resetTokenRepository) CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*db.ResetToken, error) {
	expiryDate, err := time.Parse("2006-01-02 15:04:05", req.ExpiredAt)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	res, err := r.db.CreateResetToken(ctx, db.CreateResetTokenParams{
		UserID:     int64(req.UserID),
		Token:      req.ResetToken,
		ExpiryDate: expiryDate,
	})
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}

// CreateResetTokenInTx persists the reset token inside the given database
// transaction so the caller can commit the token write and its outbox event
// atomically (Phase 6).
func (r *resetTokenRepository) CreateResetTokenInTx(ctx context.Context, tx pgx.Tx, req *requests.CreateResetTokenRequest) (*db.ResetToken, error) {
	expiryDate, err := time.Parse("2006-01-02 15:04:05", req.ExpiredAt)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	res, err := r.db.WithTx(tx).CreateResetToken(ctx, db.CreateResetTokenParams{
		UserID:     int64(req.UserID),
		Token:      req.ResetToken,
		ExpiryDate: expiryDate,
	})
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}

func (r *resetTokenRepository) DeleteResetToken(ctx context.Context, user_id int) error {
	err := r.db.DeleteResetToken(ctx, int64(user_id))
	if err != nil {
		return sharedErrors.ErrInternal.WithInternal(err)
	}
	return nil
}
