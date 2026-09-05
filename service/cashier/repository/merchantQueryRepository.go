package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/microservice-point-of-sale-cashier/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/merchant_errors"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/jackc/pgx/v5/pgtype"
)

type merchantQueryRepository struct {
	client pb.MerchantServiceClient
}

func NewMerchantQueryRepository(client pb.MerchantServiceClient) MerchantQueryRepository {
	return &merchantQueryRepository{
		client: client,
	}
}

func (r *merchantQueryRepository) FindById(ctx context.Context, id int) (*db.Merchant, error) {
	res, err := r.client.FindById(ctx, &pb.FindByIdMerchantRequest{
		Id: int32(id),
	})
	if err != nil || res == nil || res.Data == nil {
		return nil, merchant_errors.ErrFindById
	}

	var description, address, contactEmail, contactPhone *string
	if res.Data.Description != "" {
		description = &res.Data.Description
	}
	if res.Data.Address != "" {
		address = &res.Data.Address
	}
	if res.Data.ContactEmail != "" {
		contactEmail = &res.Data.ContactEmail
	}
	if res.Data.ContactPhone != "" {
		contactPhone = &res.Data.ContactPhone
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

	return &db.Merchant{
		MerchantID:   res.Data.Id,
		UserID:       res.Data.UserId,
		Name:         res.Data.Name,
		Description:  description,
		Address:      address,
		ContactEmail: contactEmail,
		ContactPhone: contactPhone,
		Status:       res.Data.Status,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}
