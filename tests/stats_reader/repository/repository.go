package repository

import "context"

// Result models mirror the ClickHouse aggregations in clickhouse.go. Counts and
// amounts are int64: ClickHouse UInt64/Int64 columns scan into them directly.

// ClickHouse returns UInt64 for count()/id columns and Int64 for sum() over
// Int64 columns; the structs mirror those exact types (clickhouse-go v2 rejects
// implicit UInt64→int64 conversion).
type MonthlyRevenue struct {
	Year         string
	Month        string
	TotalRevenue int64
	OrderCount   uint64
}

type YearlyRevenue struct {
	Year         string
	TotalRevenue int64
	OrderCount   uint64
}

type CashierMonthlyRevenue struct {
	Year         string
	Month        string
	CashierID    uint64
	TotalRevenue int64
	OrderCount   uint64
}

type ProductMonthlySold struct {
	Month     string
	ProductID uint64
	Quantity  uint64
	Subtotal  int64
}

type CategoryMonthlySold struct {
	Month      string
	CategoryID uint64
	Quantity   uint64
	Subtotal   int64
}

type TransactionMonthlySuccess struct {
	Month       string
	TotalCount  uint64
	TotalAmount int64
}

type CashierMonthlyOrders struct {
	Month       string
	CashierID   uint64
	OrderCount  uint64
	TotalAmount int64
}

// Repository serves the stats gRPC contracts (F4 §7.4) straight from
// ClickHouse. Queries never join PostgreSQL: order/cashier stats read
// order_daily, product/category stats read order_item_daily, and transaction
// stats read transaction_daily. All reads use FINAL so ReplacingMergeTree
// dedupes at-least-once redeliveries before aggregating.
type Repository interface {
	GetMonthlyTotalRevenue(ctx context.Context, year, month int) ([]MonthlyRevenue, error)
	GetYearlyTotalRevenue(ctx context.Context, year int) ([]YearlyRevenue, error)
	GetCashierMonthlyRevenue(ctx context.Context, cashierID int) ([]CashierMonthlyRevenue, error)
	GetProductMonthlySold(ctx context.Context, year, month int) ([]ProductMonthlySold, error)
	GetCategoryMonthlySold(ctx context.Context, year, month int) ([]CategoryMonthlySold, error)
	GetTransactionMonthlySuccess(ctx context.Context, year, month int) ([]TransactionMonthlySuccess, error)
	GetCashierMonthlyOrders(ctx context.Context, cashierID int) ([]CashierMonthlyOrders, error)
}
