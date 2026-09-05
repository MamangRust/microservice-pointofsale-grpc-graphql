package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-auth/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-point-of-sale-shared/errors"
	"github.com/jackc/pgx/v5"
)

type userRoleRepository struct {
	db *db.Queries
}

func NewUserRoleRepository(db *db.Queries) *userRoleRepository {
	return &userRoleRepository{
		db: db,
	}
}

func (r *userRoleRepository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*db.UserRole, error) {
	res, err := r.db.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	})
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}

// AssignRoleToUserInTx assigns the role inside the given database transaction so
// the caller can commit the role write and its outbox event atomically (Phase 6).
func (r *userRoleRepository) AssignRoleToUserInTx(ctx context.Context, tx pgx.Tx, req *requests.CreateUserRoleRequest) (*db.UserRole, error) {
	res, err := r.db.WithTx(tx).AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	})
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}

func (r *userRoleRepository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	err := r.db.RemoveRoleFromUser(ctx, db.RemoveRoleFromUserParams{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	})
	if err != nil {
		return sharedErrors.ErrInternal.WithInternal(err)
	}
	return nil
}
