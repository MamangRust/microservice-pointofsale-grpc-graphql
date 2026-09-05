package seeder

import (
	transactiondb "github.com/MamangRust/microservice-point-of-sale-transacton/database/schema"
	merchantdb "github.com/MamangRust/microservice-point-of-sale-merchant/database/schema"
	orderdb "github.com/MamangRust/microservice-point-of-sale-order/database/schema"
	"context"
	"math/rand"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"

	"go.uber.org/zap"
)

type transactionSeeder struct {
	transactiondb *transactiondb.Queries
	merchantdb *merchantdb.Queries
	orderdb *orderdb.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewTransactionSeeder(transactiondb *transactiondb.Queries, merchantdb *merchantdb.Queries, orderdb *orderdb.Queries, ctx context.Context, logger logger.LoggerInterface) *transactionSeeder {
	return &transactionSeeder{
		transactiondb: transactiondb,
		merchantdb: merchantdb,
		orderdb: orderdb,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *transactionSeeder) Seed() error {
	orders, err := r.orderdb.GetOrders(r.ctx, orderdb.GetOrdersParams{
		Column1: "",
		Limit:   20,
		Offset:  0,
	})

	if err != nil {
		r.logger.Error("Failed to get transactions:", zap.Any("error", err))
		return err
	}

	merchants, err := r.merchantdb.GetMerchants(r.ctx, merchantdb.GetMerchantsParams{
		Column1: "",
		Limit:   20,
		Offset:  0,
	})

	if err != nil {
		r.logger.Error("Failed to get transactions:", zap.Any("error", err))
		return err
	}

	for i := 0; i < 10; i++ {
		selectedMerchantId := merchants[rand.Intn(len(merchants))]
		selectedOrderId := orders[rand.Intn(len(orders))]

		var paymentMethod string
		var amount, changeAmount float64
		var paymentStatus string

		paymentMethod = "Credit Card"
		amount = float64(100 + i)
		changeAmount = float64(5 + i)
		paymentStatus = "Completed"

		_, err := r.transactiondb.CreateTransaction(r.ctx, transactiondb.CreateTransactionParams{
			OrderID:       selectedOrderId.OrderID,
			PaymentMethod: paymentMethod,
			Amount:        int32(amount),
			ChangeAmount:  ptrInt32(int32(changeAmount)),
			PaymentStatus: paymentStatus,
			MerchantID:    selectedMerchantId.MerchantID,
		})
		if err != nil {
			r.logger.Error("Failed to create transaction:", zap.Any("error", err))
			return err
		}
	}

	r.logger.Info("Successfully seeded 10 transactions.", zap.Int("count", 10))
	return nil
}
