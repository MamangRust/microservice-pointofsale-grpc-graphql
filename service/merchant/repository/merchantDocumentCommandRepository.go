package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-merchant/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-point-of-sale-shared/errors"
	"github.com/jackc/pgx/v5"
)

type merchantDocumentCommandRepository struct {
	db *db.Queries
}

func NewMerchantDocumentCommandRepository(db *db.Queries) MerchantDocumentCommandRepository {
	return &merchantDocumentCommandRepository{
		db: db,
	}
}

func (r *merchantDocumentCommandRepository) CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*db.MerchantDocument, error) {
	emptyNote := ""
	req := db.CreateMerchantDocumentParams{
		MerchantID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       "pending",
		Note:         &emptyNote,
	}

	res, err := r.db.CreateMerchantDocument(ctx, req)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

// CreateMerchantDocumentInTx persists the document inside the given database
// transaction so the caller can commit the business write and its outbox event
// atomically (Phase 6 — transactional outbox).
func (r *merchantDocumentCommandRepository) CreateMerchantDocumentInTx(ctx context.Context, tx pgx.Tx, request *requests.CreateMerchantDocumentRequest) (*db.MerchantDocument, error) {
	emptyNote := ""
	req := db.CreateMerchantDocumentParams{
		MerchantID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       "pending",
		Note:         &emptyNote,
	}

	res, err := r.db.WithTx(tx).CreateMerchantDocument(ctx, req)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

// UpdateMerchantDocumentStatusInTx updates the document status inside the given
// database transaction so the caller can commit the business write and its
// outbox event atomically (Phase 6 — transactional outbox).
func (r *merchantDocumentCommandRepository) UpdateMerchantDocumentStatusInTx(ctx context.Context, tx pgx.Tx, request *requests.UpdateMerchantDocumentStatusRequest) (*db.MerchantDocument, error) {
	if request.DocumentID == nil {
		return nil, sharedErrors.ErrBadRequest.WithMessage("document id is required")
	}

	req := db.UpdateMerchantDocumentStatusParams{
		DocumentID: int32(*request.DocumentID),
		Status:     request.Status,
		Note:       &request.Note,
	}

	res, err := r.db.WithTx(tx).UpdateMerchantDocumentStatus(ctx, req)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantDocumentCommandRepository) UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*db.MerchantDocument, error) {
	if request.DocumentID == nil {
		return nil, sharedErrors.ErrBadRequest.WithMessage("document id is required")
	}

	req := db.UpdateMerchantDocumentParams{
		DocumentID:   int32(*request.DocumentID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       request.Status,
		Note:         &request.Note,
	}

	res, err := r.db.UpdateMerchantDocument(ctx, req)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantDocumentCommandRepository) UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*db.MerchantDocument, error) {
	if request.DocumentID == nil {
		return nil, sharedErrors.ErrBadRequest.WithMessage("document id is required")
	}

	req := db.UpdateMerchantDocumentStatusParams{
		DocumentID: int32(*request.DocumentID),
		Status:     request.Status,
		Note:       &request.Note,
	}

	res, err := r.db.UpdateMerchantDocumentStatus(ctx, req)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantDocumentCommandRepository) TrashedMerchantDocument(ctx context.Context, documentID int) (*db.MerchantDocument, error) {
	res, err := r.db.TrashMerchantDocument(ctx, int32(documentID))
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantDocumentCommandRepository) RestoreMerchantDocument(ctx context.Context, documentID int) (*db.MerchantDocument, error) {
	res, err := r.db.RestoreMerchantDocument(ctx, int32(documentID))
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantDocumentCommandRepository) DeleteMerchantDocumentPermanent(ctx context.Context, documentID int) (bool, error) {
	err := r.db.DeleteMerchantDocumentPermanently(ctx, int32(documentID))
	if err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}

	return true, nil
}

func (r *merchantDocumentCommandRepository) RestoreAllMerchantDocument(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchantDocuments(ctx)
	if err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}

	return true, nil
}

func (r *merchantDocumentCommandRepository) DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchantDocuments(ctx)
	if err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}

	return true, nil
}
