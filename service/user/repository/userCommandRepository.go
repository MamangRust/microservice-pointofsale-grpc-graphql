package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-user/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/user_errors"
)

type userCommandRepository struct {
	db *db.Queries
}

func NewUserCommandRepository(db *db.Queries) *userCommandRepository {
	return &userCommandRepository{
		db: db,
	}
}

func (r *userCommandRepository) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*db.User, error) {
	req := db.CreateUserParams{
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}

	user, err := r.db.CreateUser(ctx, req)
	if err != nil {
		return nil, user_errors.ErrCreateUser
	}

	return user, nil
}

func (r *userCommandRepository) UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*db.User, error) {
	req := db.UpdateUserParams{
		UserID:    int32(*request.UserID),
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}

	res, err := r.db.UpdateUser(ctx, req)
	if err != nil {
		return nil, user_errors.ErrUpdateUser
	}

	return res, nil
}

func (r *userCommandRepository) TrashedUser(ctx context.Context, user_id int) (*db.User, error) {
	res, err := r.db.TrashUser(ctx, int32(user_id))
	if err != nil {
		return nil, user_errors.ErrTrashedUser
	}

	return res, nil
}

func (r *userCommandRepository) RestoreUser(ctx context.Context, user_id int) (*db.User, error) {
	res, err := r.db.RestoreUser(ctx, int32(user_id))
	if err != nil {
		return nil, user_errors.ErrRestoreUser
	}

	return res, nil
}

func (r *userCommandRepository) DeleteUserPermanent(ctx context.Context, user_id int) (bool, error) {
	err := r.db.DeleteUserPermanently(ctx, int32(user_id))
	if err != nil {
		return false, user_errors.ErrDeleteUserPermanent
	}

	return true, nil
}

func (r *userCommandRepository) RestoreAllUser(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllUsers(ctx)
	if err != nil {
		return false, user_errors.ErrRestoreAllUsers
	}

	return true, nil
}

func (r *userCommandRepository) DeleteAllUserPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentUsers(ctx)
	if err != nil {
		return false, user_errors.ErrDeleteAllUsers
	}
	return true, nil
}
