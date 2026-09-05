package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/microservice-point-of-sale-cashier/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/user_errors"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/jackc/pgx/v5/pgtype"
)

type userQueryRepository struct {
	client pb.UserServiceClient
}

func NewUserQueryRepository(client pb.UserServiceClient) UserQueryRepository {
	return &userQueryRepository{
		client: client,
	}
}

func (r *userQueryRepository) FindById(ctx context.Context, id int) (*db.User, error) {
	res, err := r.client.FindById(ctx, &pb.FindByIdUserRequest{
		Id: int32(id),
	})
	if err != nil || res == nil || res.Data == nil {
		return nil, user_errors.ErrUserNotFound
	}

	var createdAt, updatedAt pgtype.Timestamp
	if res.Data.CreatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", res.Data.CreatedAt); err == nil {
			createdAt = pgtype.Timestamp{Time: t, Valid: true}
		} else if t, err = time.Parse(time.RFC3339, res.Data.CreatedAt); err == nil {
			createdAt = pgtype.Timestamp{Time: t, Valid: true}
		}
	}
	if res.Data.UpdatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", res.Data.UpdatedAt); err == nil {
			updatedAt = pgtype.Timestamp{Time: t, Valid: true}
		} else if t, err = time.Parse(time.RFC3339, res.Data.UpdatedAt); err == nil {
			updatedAt = pgtype.Timestamp{Time: t, Valid: true}
		}
	}

	return &db.User{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
