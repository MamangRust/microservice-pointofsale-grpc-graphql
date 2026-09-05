package main

import (
	"context"
	"log"

	cashierdb "github.com/MamangRust/microservice-point-of-sale-cashier/database/schema"
	categorydb "github.com/MamangRust/microservice-point-of-sale-category/database/schema"
	merchantdb "github.com/MamangRust/microservice-point-of-sale-merchant/database/schema"
	orderdb "github.com/MamangRust/microservice-point-of-sale-order/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-pkg/database"
	"github.com/MamangRust/microservice-point-of-sale-pkg/database/seeder"
	"github.com/MamangRust/microservice-point-of-sale-pkg/dotenv"
	"github.com/MamangRust/microservice-point-of-sale-pkg/hash"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	productdb "github.com/MamangRust/microservice-point-of-sale-product/database/schema"
	roledb "github.com/MamangRust/microservice-point-of-sale-role/database/schema"
	transactiondb "github.com/MamangRust/microservice-point-of-sale-transacton/database/schema"
	userdb "github.com/MamangRust/microservice-point-of-sale-user/database/schema"
)

func main() {
	if err := dotenv.Viper(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	l, err := logger.NewLogger("seeder", nil)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	pool, err := database.NewClient(l)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer pool.Close()

	ctx := context.Background()

	s := seeder.NewSeeder(seeder.Deps{
		User:        userdb.New(pool),
		Role:        roledb.New(pool),
		Merchant:    merchantdb.New(pool),
		Cashier:     cashierdb.New(pool),
		Category:    categorydb.New(pool),
		Product:     productdb.New(pool),
		Order:       orderdb.New(pool),
		Transaction: transactiondb.New(pool),
		Ctx:         ctx,
		Logger:      l,
		Hash:        hash.NewHashingPassword(),
	})

	if err := s.Run(); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	l.Info("Seeding completed successfully.")
}
