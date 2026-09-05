package order_test

import (
	"context"
	orderdb "github.com/MamangRust/microservice-point-of-sale-order/database/schema"
	"testing"

	"github.com/MamangRust/microservice-point-of-sale-order/repository"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	tests "github.com/MamangRust/microservice-point-of-sale-test"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	repo       *repository.Repositories
	orderID    int
	merchantID int
	cashierID  int
}

func (s *OrderRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)

	// Create placeholder gRPC connections (gRPC repos not called in DB-only tests)
	cashierLis, _ := net.Listen("tcp", "localhost:0")
	merchantLis, _ := net.Listen("tcp", "localhost:0")
	productLis, _ := net.Listen("tcp", "localhost:0")
	orderItemLis, _ := net.Listen("tcp", "localhost:0")

	cashierConn, _ := grpc.NewClient(cashierLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	merchantConn, _ := grpc.NewClient(merchantLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	productConn, _ := grpc.NewClient(productLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	orderItemConn, _ := grpc.NewClient(orderItemLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	queries := orderdb.New(pool)
	s.repo = repository.NewRepositories(
		queries,
		pb.NewCashierServiceClient(cashierConn),
		pb.NewMerchantServiceClient(merchantConn),
		pb.NewProductServiceClient(productConn),
		pb.NewOrderItemServiceClient(orderItemConn),
	)

	// Seed a merchant and cashier directly for order tests
	var userID int
	err = pool.QueryRow(s.ts.Ctx,
		`INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ($1, $2, $3, $4, 'test-verify', true) RETURNING user_id`,
		"Order", "Repo", "order.repo@example.com", "password123",
	).Scan(&userID)
	s.Require().NoError(err)

	err = pool.QueryRow(s.ts.Ctx,
		`INSERT INTO merchants (user_id, name, description, address, contact_email, contact_phone, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING merchant_id`,
		userID, "Order Test Merchant", "Desc", "Addr", "o@example.com", "123", "active",
	).Scan(&s.merchantID)
	s.Require().NoError(err)

	// Seed a cashier (orders.cashier_id references cashiers.cashier_id)
	err = pool.QueryRow(s.ts.Ctx,
		`INSERT INTO cashiers (merchant_id, user_id, name) VALUES ($1, $2, 'Order Test Cashier') RETURNING cashier_id`,
		s.merchantID, userID,
	).Scan(&s.cashierID)
	s.Require().NoError(err)

	// Seed an order for use in tests
	err = pool.QueryRow(s.ts.Ctx,
		`INSERT INTO orders (merchant_id, cashier_id, total_price) VALUES ($1, $2, $3) RETURNING order_id`,
		s.merchantID, s.cashierID, 5000,
	).Scan(&s.orderID)
	s.Require().NoError(err)
}

func (s *OrderRepositoryTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *OrderRepositoryTestSuite) Test1_FindById() {
	s.Require().NotZero(s.orderID)
	ctx := context.Background()

	found, err := s.repo.OrderQuery.FindById(ctx, s.orderID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(int32(s.orderID), found.OrderID)
}

func (s *OrderRepositoryTestSuite) Test2_CreateOrder() {
	ctx := context.Background()

	req := &requests.CreateOrderRecordRequest{
		MerchantID: s.merchantID,
		CashierID:  s.cashierID,
		TotalPrice: 10000,
	}

	res, err := s.repo.OrderCommand.CreateOrder(ctx, req)
	s.NoError(err)
	s.NotNil(res)
}

func (s *OrderRepositoryTestSuite) Test3_UpdateOrder() {
	s.Require().NotZero(s.orderID)
	ctx := context.Background()

	req := &requests.UpdateOrderRecordRequest{
		OrderID:    s.orderID,
		TotalPrice: 15000,
	}

	res, err := s.repo.OrderCommand.UpdateOrder(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(int64(15000), res.TotalPrice)
}

func (s *OrderRepositoryTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.orderID)
	ctx := context.Background()

	// Trash
	trashed, err := s.repo.OrderCommand.TrashedOrder(ctx, s.orderID)
	s.NoError(err)
	s.NotNil(trashed)

	// Restore
	restored, err := s.repo.OrderCommand.RestoreOrder(ctx, s.orderID)
	s.NoError(err)
	s.NotNil(restored)

	// Verify restored
	found, err := s.repo.OrderQuery.FindById(ctx, s.orderID)
	s.NoError(err)
	s.NotNil(found)
}

func (s *OrderRepositoryTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.orderID)
	ctx := context.Background()

	// Must be trashed first for permanent delete
	_, err := s.repo.OrderCommand.TrashedOrder(ctx, s.orderID)
	s.NoError(err)

	success, err := s.repo.OrderCommand.DeleteOrderPermanent(ctx, s.orderID)
	s.NoError(err)
	s.True(success)
}

func TestOrderRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderRepositoryTestSuite))
}
