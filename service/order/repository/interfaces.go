package repository

import (
	"context"

	db "github.com/MamangRust/microservice-point-of-sale-order/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
)

type OrderQueryRepository interface {
	FindAllOrders(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersTrashedRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllOrderMerchant) ([]*db.GetOrdersByMerchantRow, *int, error)
	FindById(ctx context.Context, orderID int) (*db.Order, error)
	FindByTrashedId(ctx context.Context, orderID int) (*db.Order, error)
}

type OrderCommandRepository interface {
	DeleteOrder(ctx context.Context, orderID int) error
	CreateOrder(ctx context.Context, request *requests.CreateOrderRecordRequest) (*db.Order, error)
	UpdateOrder(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*db.Order, error)
	FindAllTrashed(ctx context.Context) ([]*db.Order, error)
	TrashedOrder(ctx context.Context, orderID int) (*db.Order, error)
	RestoreOrder(ctx context.Context, orderID int) (*db.Order, error)
	DeleteOrderPermanent(ctx context.Context, orderID int) (bool, error)
	DeleteAllOrderPermanent(ctx context.Context) (bool, error)
}

type CashierQueryRepository interface {
	FindById(ctx context.Context, cashierID int) (*db.Cashier, error)
}

type MerchantQueryRepository interface {
	FindById(ctx context.Context, merchantID int) (*db.Merchant, error)
}

type ProductQueryRepository interface {
	FindById(ctx context.Context, product_id int) (*db.Product, error)
}

type ProductCommandRepository interface {
	DecrementProductCountStock(ctx context.Context, productID int, quantity int) (*db.Product, error)
	IncrementProductCountStock(ctx context.Context, productID int, quantity int) (*db.Product, error)
}

type OrderItemQueryRepository interface {
	FindOrderItemByOrder(ctx context.Context, orderID int) ([]*db.OrderItem, error)
	CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error)
}

type OrderItemCommandRepository interface {
	DeleteOrderItem(ctx context.Context, orderID int) error
	CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*db.OrderItem, error)
	UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*db.OrderItem, error)
}
