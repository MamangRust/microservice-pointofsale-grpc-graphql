package repository

import (
	db "github.com/MamangRust/microservice-point-of-sale-order-item/database/schema"
)

type Repositories struct {
	OrderItemQuery OrderItemQueryRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	return &Repositories{
		OrderItemQuery: NewOrderItemQueryRepository(DB),
	}
}
