package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-merchant/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-point-of-sale-shared/errors"
	"github.com/jackc/pgx/v5"
)

type merchantCommandRepository struct {
	db *db.Queries
}

func NewMerchantCommandRepository(db *db.Queries) MerchantCommandRepository {
	return &merchantCommandRepository{
		db: db,
	}
}

func (r *merchantCommandRepository) CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*db.Merchant, error) {
	req := db.CreateMerchantParams{
		UserID:       int32(request.UserID),
		Name:         request.Name,
		Description:  &request.Description,
		Address:      &request.Address,
		ContactEmail: &request.ContactEmail,
		ContactPhone: &request.ContactPhone,
		Status:       "inactive",
	}

	merchant, err := r.db.CreateMerchant(ctx, req)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return merchant, nil
}

// CreateMerchantInTx persists the merchant inside the given database transaction
// so the caller can commit the business write and its outbox event atomically
// (Phase 6 — transactional outbox).
func (r *merchantCommandRepository) CreateMerchantInTx(ctx context.Context, tx pgx.Tx, request *requests.CreateMerchantRequest) (*db.Merchant, error) {
	req := db.CreateMerchantParams{
		UserID:       int32(request.UserID),
		Name:         request.Name,
		Description:  &request.Description,
		Address:      &request.Address,
		ContactEmail: &request.ContactEmail,
		ContactPhone: &request.ContactPhone,
		Status:       "inactive",
	}

	merchant, err := r.db.WithTx(tx).CreateMerchant(ctx, req)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return merchant, nil
}

// UpdateMerchantStatusInTx updates the merchant status inside the given database
// transaction so the caller can commit the business write and its outbox event
// atomically (Phase 6 — transactional outbox).
func (r *merchantCommandRepository) UpdateMerchantStatusInTx(ctx context.Context, tx pgx.Tx, request *requests.UpdateMerchantStatusRequest) (*db.Merchant, error) {
	req := db.UpdateMerchantStatusParams{
		MerchantID: int32(*request.MerchantID),
		Status:     request.Status,
	}

	res, err := r.db.WithTx(tx).UpdateMerchantStatus(ctx, req)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantCommandRepository) UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*db.Merchant, error) {
	req := db.UpdateMerchantParams{
		MerchantID:   int32(*request.MerchantID),
		Name:         request.Name,
		Description:  &request.Description,
		Address:      &request.Address,
		ContactEmail: &request.ContactEmail,
		ContactPhone: &request.ContactPhone,
		Status:       request.Status,
	}

	res, err := r.db.UpdateMerchant(ctx, req)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantCommandRepository) UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*db.Merchant, error) {
	req := db.UpdateMerchantStatusParams{
		MerchantID: int32(*request.MerchantID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateMerchantStatus(ctx, req)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantCommandRepository) TrashedMerchant(ctx context.Context, merchant_id int) (*db.Merchant, error) {
	res, err := r.db.TrashMerchant(ctx, int32(merchant_id))
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantCommandRepository) RestoreMerchant(ctx context.Context, merchant_id int) (*db.Merchant, error) {
	res, err := r.db.RestoreMerchant(ctx, int32(merchant_id))
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantCommandRepository) DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantPermanently(ctx, int32(merchant_id))
	if err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}

	return true, nil
}

func (r *merchantCommandRepository) RestoreAllMerchant(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchants(ctx)
	if err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}
	return true, nil
}

func (r *merchantCommandRepository) DeleteAllMerchantPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(ctx)
	if err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}
	return true, nil
}
