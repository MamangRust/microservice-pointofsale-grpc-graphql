package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-transacton/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/jackc/pgx/v5"
)

type CashierQueryRepository interface {
	FindById(ctx context.Context, id int) (*db.Cashier, error)
}

type MerchantQueryRepository interface {
	FindById(ctx context.Context, id int) (*db.Merchant, error)
}

type OrderItemQueryRepository interface {
	FindOrderItemByOrder(ctx context.Context, order_id int) ([]*db.OrderItem, error)
}

type OrderQueryRepository interface {
	FindById(ctx context.Context, id int) (*db.Order, error)
}

type TransactionQueryRepository interface {
	FindAllTransactions(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsTrashedRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*db.GetTransactionByMerchantRow, *int, error)
	FindById(ctx context.Context, transaction_id int) (*db.Transaction, error)
	FindByOrderId(ctx context.Context, order_id int) (*db.Transaction, error)
}

type TransactionCommandRepository interface {
	CreateTransaction(ctx context.Context, request *requests.CreateTransactionRequest) (*db.Transaction, error)
	CreateTransactionInTx(ctx context.Context, tx pgx.Tx, request *requests.CreateTransactionRequest) (*db.Transaction, error)
	UpdateTransaction(ctx context.Context, request *requests.UpdateTransactionRequest) (*db.Transaction, error)
	TrashTransaction(ctx context.Context, transaction_id int) (*db.Transaction, error)
	RestoreTransaction(ctx context.Context, transaction_id int) (*db.Transaction, error)
	DeleteTransactionPermanently(ctx context.Context, transaction_id int) (bool, error)
	RestoreAllTransactions(ctx context.Context) (bool, error)
	DeleteAllTransactionPermanent(ctx context.Context) (bool, error)
}
