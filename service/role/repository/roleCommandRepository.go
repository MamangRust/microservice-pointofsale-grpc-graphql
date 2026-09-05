package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-role/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
)

type roleCommandRepository struct {
	db *db.Queries
}

func NewRoleCommandRepository(db *db.Queries) *roleCommandRepository {
	return &roleCommandRepository{
		db: db,
	}
}

func (r *roleCommandRepository) CreateRole(ctx context.Context, req *requests.CreateRoleRequest) (*db.Role, error) {
	res, err := r.db.CreateRole(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *roleCommandRepository) UpdateRole(ctx context.Context, req *requests.UpdateRoleRequest) (*db.Role, error) {
	res, err := r.db.UpdateRole(ctx, db.UpdateRoleParams{
		RoleID:   int32(*req.ID),
		RoleName: req.Name,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *roleCommandRepository) TrashedRole(ctx context.Context, id int) (*db.Role, error) {
	res, err := r.db.TrashRole(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *roleCommandRepository) RestoreRole(ctx context.Context, id int) (*db.Role, error) {
	res, err := r.db.RestoreRole(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *roleCommandRepository) DeleteRolePermanent(ctx context.Context, roleID int) (bool, error) {
	err := r.db.DeletePermanentRole(ctx, int32(roleID))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *roleCommandRepository) RestoreAllRole(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllRoles(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *roleCommandRepository) DeleteAllRolePermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentRoles(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}
