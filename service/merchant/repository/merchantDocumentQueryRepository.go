package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-merchant/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-point-of-sale-shared/errors"
)

type merchantDocumentQueryRepository struct {
	db *db.Queries
}

func NewMerchantDocumentQueryRepository(db *db.Queries) MerchantDocumentQueryRepository {
	return &merchantDocumentQueryRepository{
		db: db,
	}
}

func (r *merchantDocumentQueryRepository) FindAllDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetMerchantDocumentsRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	params := db.GetMerchantDocumentsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	docs, err := r.db.GetMerchantDocuments(ctx, params)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var totalCount int
	if len(docs) > 0 {
		totalCount = int(docs[0].TotalCount)
	}

	return docs, &totalCount, nil
}

func (r *merchantDocumentQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetActiveMerchantDocumentsRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	params := db.GetActiveMerchantDocumentsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	docs, err := r.db.GetActiveMerchantDocuments(ctx, params)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var totalCount int
	if len(docs) > 0 {
		totalCount = int(docs[0].TotalCount)
	}

	return docs, &totalCount, nil
}

func (r *merchantDocumentQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetTrashedMerchantDocumentsRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	params := db.GetTrashedMerchantDocumentsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	docs, err := r.db.GetTrashedMerchantDocuments(ctx, params)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var totalCount int
	if len(docs) > 0 {
		totalCount = int(docs[0].TotalCount)
	}

	return docs, &totalCount, nil
}

func (r *merchantDocumentQueryRepository) FindById(ctx context.Context, id int) (*db.MerchantDocument, error) {
	doc, err := r.db.GetMerchantDocument(ctx, int32(id))
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return doc, nil
}
