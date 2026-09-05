package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-transacton/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/transaction_errors"
)

type transactionQueryRepository struct {
	db *db.Queries
}

func NewTransactionQueryRepository(db *db.Queries) *transactionQueryRepository {
	return &transactionQueryRepository{
		db: db,
	}
}

func (r *transactionQueryRepository) FindAllTransactions(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactions(ctx, reqDb)
	if err != nil {
		return nil, nil, transaction_errors.ErrFindAllTransactions
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *transactionQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactionsActive(ctx, reqDb)
	if err != nil {
		return nil, nil, transaction_errors.ErrFindByActive
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *transactionQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactionsTrashed(ctx, reqDb)
	if err != nil {
		return nil, nil, transaction_errors.ErrFindByTrashed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *transactionQueryRepository) FindByMerchant(
	ctx context.Context,
	req *requests.FindAllTransactionByMerchant,
) ([]*db.GetTransactionByMerchantRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionByMerchantParams{
		Column1: req.Search,
		Column2: int32(req.MerchantID),
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactionByMerchant(ctx, reqDb)
	if err != nil {
		return nil, nil, transaction_errors.ErrFindByMerchant
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return res, &totalCount, nil
}

func (r *transactionQueryRepository) FindById(ctx context.Context, transaction_id int) (*db.Transaction, error) {
	res, err := r.db.GetTransactionByID(ctx, int32(transaction_id))
	if err != nil {
		return nil, transaction_errors.ErrFindById
	}

	return res, nil
}

func (r *transactionQueryRepository) FindByOrderId(ctx context.Context, order_id int) (*db.Transaction, error) {
	res, err := r.db.GetTransactionByOrderID(ctx, int32(order_id))
	if err != nil {
		return nil, transaction_errors.ErrFindByOrderId
	}

	return res, nil
}
