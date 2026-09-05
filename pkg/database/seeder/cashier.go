package seeder

import (
	cashierdb "github.com/MamangRust/microservice-point-of-sale-cashier/database/schema"
	merchantdb "github.com/MamangRust/microservice-point-of-sale-merchant/database/schema"
	userdb "github.com/MamangRust/microservice-point-of-sale-user/database/schema"
	"context"
	"fmt"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"

	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

type cashierSeeder struct {
	cashierdb *cashierdb.Queries
	merchantdb *merchantdb.Queries
	userdb *userdb.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewCashierSeeder(cashierdb *cashierdb.Queries, merchantdb *merchantdb.Queries, userdb *userdb.Queries, ctx context.Context, logger logger.LoggerInterface) *cashierSeeder {
	return &cashierSeeder{
		cashierdb: cashierdb,
		merchantdb: merchantdb,
		userdb: userdb,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *cashierSeeder) Seed() error {
	merchants, err := r.merchantdb.GetMerchants(r.ctx, merchantdb.GetMerchantsParams{
		Column1: "",
		Limit:   int32(20),
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("Failed to fetch merchants:", zap.Any("err", err))
		return err
	}

	users, err := r.userdb.GetUsers(r.ctx, userdb.GetUsersParams{
		Column1: "",
		Limit:   int32(20),
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("Failed to fetch users:", zap.Any("error", err))
		return err
	}

	if len(merchants) == 0 || len(users) == 0 {
		r.logger.Error("Merchants or Users not found. Seed operation aborted.")
		return fmt.Errorf("no merchants or users found")
	}

	for i := 1; i <= 10; i++ {
		merchant := merchants[rand.Intn(len(merchants))]
		user := users[rand.Intn(len(users))]

		cashierName := fmt.Sprintf("Cashier %d", i)
		_, err := r.cashierdb.CreateCashier(r.ctx, cashierdb.CreateCashierParams{
			MerchantID: merchant.MerchantID,
			UserID:     user.UserID,
			Name:       cashierName,
		})
		if err != nil {
			r.logger.Error("Failed to create cashier:", zap.Any("error", err))
			return err
		}

	}

	r.logger.Info("Cashier seeding completed successfully.", zap.Int("count", 10))
	return nil
}
