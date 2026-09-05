package repository

import (
	"context"
	"errors"

	db "github.com/MamangRust/microservice-point-of-sale-user/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/role_errors"
	"github.com/jackc/pgx/v5"
)

type roleRepository struct {
	db *db.Queries
}

func NewRoleRepository(db *db.Queries) *roleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*db.Role, error) {
	res, err := r.db.GetRoleByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, role_errors.ErrRoleNotFound
		}
		return nil, role_errors.ErrRoleNotFound
	}
	return res, nil
}
