package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-auth/database/schema"
	sharedErrors "github.com/MamangRust/microservice-point-of-sale-shared/errors"
)

type roleRepository struct {
	db *db.Queries
}

func NewRoleRepository(db *db.Queries) *roleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) FindById(ctx context.Context, id int) (*db.Role, error) {
	res, err := r.db.GetRole(ctx, int32(id))
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*db.Role, error) {
	res, err := r.db.GetRoleByName(ctx, name)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}
