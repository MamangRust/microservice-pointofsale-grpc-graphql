package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-transacton/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/transaction_errors"
	"github.com/jackc/pgx/v5"
)

type transactionCommandRepository struct {
	db *db.Queries
}

func NewTransactionCommandRepository(db *db.Queries) *transactionCommandRepository {
	return &transactionCommandRepository{
		db: db,
	}
}

func (r *transactionCommandRepository) CreateTransaction(ctx context.Context, request *requests.CreateTransactionRequest) (*db.Transaction, error) {
	req := db.CreateTransactionParams{
		OrderID:       int32(request.OrderID),
		MerchantID:    int32(request.MerchantID),
		PaymentMethod: request.PaymentMethod,
		Amount:        int32(request.Amount),
		PaymentStatus: *request.PaymentStatus,
	}

	transaction, err := r.db.CreateTransaction(ctx, req)
	if err != nil {
		return nil, transaction_errors.ErrCreateTransaction
	}

	return transaction, nil
}

// CreateTransactionInTx persists the transaction inside the given database
// transaction so the caller can commit the business write and its outbox event
// atomically (Phase 6 — transactional outbox).
func (r *transactionCommandRepository) CreateTransactionInTx(ctx context.Context, tx pgx.Tx, request *requests.CreateTransactionRequest) (*db.Transaction, error) {
	req := db.CreateTransactionParams{
		OrderID:       int32(request.OrderID),
		MerchantID:    int32(request.MerchantID),
		PaymentMethod: request.PaymentMethod,
		Amount:        int32(request.Amount),
		PaymentStatus: *request.PaymentStatus,
	}

	transaction, err := r.db.WithTx(tx).CreateTransaction(ctx, req)
	if err != nil {
		return nil, transaction_errors.ErrCreateTransaction
	}

	return transaction, nil
}

func (r *transactionCommandRepository) UpdateTransaction(ctx context.Context, request *requests.UpdateTransactionRequest) (*db.Transaction, error) {
	req := db.UpdateTransactionParams{
		TransactionID: int32(*request.TransactionID),
		MerchantID:    int32(request.MerchantID),
		PaymentMethod: request.PaymentMethod,
		Amount:        int32(request.Amount),
		OrderID:       int32(request.OrderID),
		PaymentStatus: *request.PaymentStatus,
	}

	res, err := r.db.UpdateTransaction(ctx, req)
	if err != nil {
		return nil, transaction_errors.ErrUpdateTransaction
	}

	return res, nil
}

func (r *transactionCommandRepository) TrashTransaction(ctx context.Context, transaction_id int) (*db.Transaction, error) {
	res, err := r.db.TrashTransaction(ctx, int32(transaction_id))
	if err != nil {
		return nil, transaction_errors.ErrTrashTransaction
	}

	return res, nil
}

func (r *transactionCommandRepository) RestoreTransaction(ctx context.Context, transaction_id int) (*db.Transaction, error) {
	res, err := r.db.RestoreTransaction(ctx, int32(transaction_id))
	if err != nil {
		return nil, transaction_errors.ErrRestoreTransaction
	}

	return res, nil
}

func (r *transactionCommandRepository) DeleteTransactionPermanently(ctx context.Context, transaction_id int) (bool, error) {
	err := r.db.DeleteTransactionPermanently(ctx, int32(transaction_id))
	if err != nil {
		return false, transaction_errors.ErrDeleteTransactionPermanently
	}

	return true, nil
}

func (r *transactionCommandRepository) RestoreAllTransactions(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllTransactions(ctx)
	if err != nil {
		return false, transaction_errors.ErrRestoreAllTransactions
	}
	return true, nil
}

func (r *transactionCommandRepository) DeleteAllTransactionPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentTransactions(ctx)
	if err != nil {
		return false, transaction_errors.ErrDeleteAllTransactionPermanent
	}
	return true, nil
}
