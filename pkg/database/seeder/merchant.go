package seeder

import (
	merchantdb "github.com/MamangRust/microservice-point-of-sale-merchant/database/schema"
	userdb "github.com/MamangRust/microservice-point-of-sale-user/database/schema"
	"context"
	"fmt"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"

	"go.uber.org/zap"
)

type merchantSeeder struct {
	merchantdb *merchantdb.Queries
	userdb *userdb.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantSeeder(merchantdb *merchantdb.Queries, userdb *userdb.Queries, ctx context.Context, logger logger.LoggerInterface) *merchantSeeder {
	return &merchantSeeder{
		merchantdb: merchantdb,
		userdb: userdb,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *merchantSeeder) Seed() error {
	users, err := r.userdb.GetUsers(r.ctx, userdb.GetUsersParams{
		Column1: "",
		Limit:   int32(20),
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("Failed to fetch merchants:", zap.Any("error", err))
		return err
	}

	for i := 1; i <= 10; i++ {
		userID := users[i%len(users)].UserID

		merchant := merchantdb.CreateMerchantParams{
			UserID:       userID,
			Name:         fmt.Sprintf("Toko %d", i),
			Description:  ptrString(fmt.Sprintf("Deskripsi untuk Toko %d", i)),
			Address:      ptrString(fmt.Sprintf("Jl. Toko %d", i)),
			ContactEmail: ptrString(fmt.Sprintf("toko%d@example.com", i)),
			ContactPhone: ptrString(fmt.Sprintf("0812345678%d", i)),
			Status:       "active",
		}

		_, err = r.merchantdb.CreateMerchant(r.ctx, merchant)
		if err != nil {
			r.logger.Error("Failed to create merchant:", zap.Error(err))
			return err
		}
	}

	r.logger.Info("Merchant seeding completed successfully.", zap.Int("count", 10))
	return nil
}
