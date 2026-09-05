package repository

import (
	"context"
	"errors"

	db "github.com/MamangRust/microservice-point-of-sale-auth/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-point-of-sale-shared/errors"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/user_errors"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	db *db.Queries
}

func NewUserRepository(db *db.Queries) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) FindById(ctx context.Context, user_id int) (*db.User, error) {
	res, err := r.db.GetUserByID(ctx, int32(user_id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*db.User, error) {
	res, err := r.db.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}

func (r *userRepository) FindByEmailAndVerify(ctx context.Context, email string) (*db.User, error) {
	res, err := r.db.GetUserByEmailAndVerified(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}

func (r *userRepository) FindByVerificationCode(ctx context.Context, verification_code string) (*db.User, error) {
	res, err := r.db.GetUserByVerificationCode(ctx, verification_code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}

func (r *userRepository) CreateUser(ctx context.Context, request *requests.RegisterRequest) (*db.User, error) {
	req := db.CreateUserParams{
		Firstname:        request.FirstName,
		Lastname:         request.LastName,
		Email:            request.Email,
		Password:         request.Password,
		VerificationCode: request.VerifiedCode,
		IsVerified:       &request.IsVerified,
	}

	user, err := r.db.CreateUser(ctx, req)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return user, nil
}

// CreateUserInTx persists the user inside the given database transaction so the
// caller can commit the user write and its outbox event atomically (Phase 6).
func (r *userRepository) CreateUserInTx(ctx context.Context, tx pgx.Tx, request *requests.RegisterRequest) (*db.User, error) {
	req := db.CreateUserParams{
		Firstname:        request.FirstName,
		Lastname:         request.LastName,
		Email:            request.Email,
		Password:         request.Password,
		VerificationCode: request.VerifiedCode,
		IsVerified:       &request.IsVerified,
	}

	user, err := r.db.WithTx(tx).CreateUser(ctx, req)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return user, nil
}

func (r *userRepository) UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*db.User, error) {
	res, err := r.db.UpdateUserIsVerified(ctx, db.UpdateUserIsVerifiedParams{
		UserID:     int32(user_id),
		IsVerified: &is_verified,
	})
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}

func (r *userRepository) UpdateUserPassword(ctx context.Context, user_id int, password string) (*db.User, error) {
	res, err := r.db.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		UserID:   int32(user_id),
		Password: password,
	})
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}
