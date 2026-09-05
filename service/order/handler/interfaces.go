package handler

import "github.com/MamangRust/microservice-point-of-sale-shared/pb"

type OrderHandleGrpc interface {
	pb.OrderServiceServer
}
