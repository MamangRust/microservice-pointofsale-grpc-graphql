package handler

import "github.com/MamangRust/microservice-point-of-sale-shared/pb"

type ProductHandleGrpc interface {
	pb.ProductServiceServer
}
