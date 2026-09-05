package repository

import (
	db "github.com/MamangRust/microservice-point-of-sale-cashier/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
)

type Repositories struct {
	UserQuery              UserQueryRepository
	MerchantQuery          MerchantQueryRepository
	CashierQuery           CashierQueryRepository
	CashierCommand         CashierCommandRepository
}

func NewRepositories(
	DB *db.Queries,
	userClient pb.UserServiceClient,
	merchantClient pb.MerchantServiceClient,
) *Repositories {
	return &Repositories{
		UserQuery:              NewUserQueryRepository(userClient),
		MerchantQuery:          NewMerchantQueryRepository(merchantClient),
		CashierQuery:           NewCashierQueryRepository(DB),
		CashierCommand:         NewCashierCommandRepository(DB),
	}
}
