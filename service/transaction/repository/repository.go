package repository

import (
	"github.com/MamangRust/microservice-point-of-sale-pkg/adapter"
	db "github.com/MamangRust/microservice-point-of-sale-transacton/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
)

type Repositories struct {
	CashierQuery                 CashierQueryRepository
	MerchantQuery                MerchantQueryRepository
	OrderQuery                   OrderQueryRepository
	OrderItemQuery               OrderItemQueryRepository
	TransactionCommandRepository TransactionCommandRepository
	TransactionQueryRepository   TransactionQueryRepository
}

// GuardOptions groups dependency guards by downstream gRPC dependency. Each
// non-nil slice is applied to the matching repository at construction (F6 §8.1
// poin 5 — deadline + circuit breaker + bulkhead per dependency).
type GuardOptions struct {
	Cashier     []adapter.GuardOption
	Merchant    []adapter.GuardOption
	Order       []adapter.GuardOption
	OrderItem   []adapter.GuardOption
}

func NewRepositories(
	DB *db.Queries,
	cashierClient pb.CashierServiceClient,
	merchantClient pb.MerchantServiceClient,
	orderClient pb.OrderServiceClient,
	orderItemClient pb.OrderItemServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}
	return &Repositories{
		CashierQuery:                 NewCashierQueryRepository(cashierClient, g.Cashier...),
		MerchantQuery:                NewMerchantQueryRepository(merchantClient, g.Merchant...),
		OrderQuery:                   NewOrderQueryRepository(orderClient, g.Order...),
		OrderItemQuery:               NewOrderItemQueryRepository(orderItemClient, g.OrderItem...),
		TransactionCommandRepository: NewTransactionCommandRepository(DB),
		TransactionQueryRepository:   NewTransactionQueryRepository(DB),
	}
}
