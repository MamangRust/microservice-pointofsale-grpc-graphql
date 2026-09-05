package order_test

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-order/database/schema"
)

// stubCashierRepo implements repository.CashierQueryRepository. In the
// failure-injection suite it is never reached (the merchant dependency fails
// first), so it returns a minimal valid row to keep the flow honest if the
// ordering ever changes.
type stubCashierRepo struct {
	cashierID int
}

func (s *stubCashierRepo) FindById(ctx context.Context, cashierID int) (*db.Cashier, error) {
	return &db.Cashier{CashierID: int32(cashierID), MerchantID: int32(s.cashierID), Name: "stub"}, nil
}

// stubProductRepo implements repository.ProductQueryRepository. Never reached
// in the failure-injection suite (merchant fails first).
type stubProductRepo struct{}

func (s *stubProductRepo) FindById(ctx context.Context, productID int) (*db.Product, error) {
	return &db.Product{ProductID: int32(productID), Price: 10000, CountInStock: 100}, nil
}
