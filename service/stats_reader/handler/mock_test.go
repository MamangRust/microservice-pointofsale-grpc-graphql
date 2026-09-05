package handler

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-stats-reader/repository"
	"go.uber.org/zap"
)

type mockRepo struct {
	monthlyRevenue     []repository.MonthlyRevenue
	yearlyRevenue      []repository.YearlyRevenue
	cashierRevenue     []repository.CashierMonthlyRevenue
	productSold        []repository.ProductMonthlySold
	categorySold       []repository.CategoryMonthlySold
	transactionSuccess []repository.TransactionMonthlySuccess
	cashierOrders      []repository.CashierMonthlyOrders
	err                error
}

func (m *mockRepo) GetMonthlyTotalRevenue(_ context.Context, _, _ int) ([]repository.MonthlyRevenue, error) {
	return m.monthlyRevenue, m.err
}

func (m *mockRepo) GetYearlyTotalRevenue(_ context.Context, _ int) ([]repository.YearlyRevenue, error) {
	return m.yearlyRevenue, m.err
}

func (m *mockRepo) GetCashierMonthlyRevenue(_ context.Context, _ int) ([]repository.CashierMonthlyRevenue, error) {
	return m.cashierRevenue, m.err
}

func (m *mockRepo) GetProductMonthlySold(_ context.Context, _, _ int) ([]repository.ProductMonthlySold, error) {
	return m.productSold, m.err
}

func (m *mockRepo) GetCategoryMonthlySold(_ context.Context, _, _ int) ([]repository.CategoryMonthlySold, error) {
	return m.categorySold, m.err
}

func (m *mockRepo) GetTransactionMonthlySuccess(_ context.Context, _, _ int) ([]repository.TransactionMonthlySuccess, error) {
	return m.transactionSuccess, m.err
}

func (m *mockRepo) GetCashierMonthlyOrders(_ context.Context, _ int) ([]repository.CashierMonthlyOrders, error) {
	return m.cashierOrders, m.err
}

// stubLogger satisfies logger.LoggerInterface
type stubLogger struct{}

func (s *stubLogger) Info(_ string, _ ...zap.Field)  {}
func (s *stubLogger) Fatal(_ string, _ ...zap.Field) {}
func (s *stubLogger) Debug(_ string, _ ...zap.Field) {}
func (s *stubLogger) Error(_ string, _ ...zap.Field) {}
func (s *stubLogger) Warn(_ string, _ ...zap.Field)  {}
