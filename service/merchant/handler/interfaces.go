package handler

import "github.com/MamangRust/microservice-point-of-sale-shared/pb"

type MerchantDocumentHandleGrpc interface {
	pb.MerchantDocumentServiceServer
}

type MerchantHandleGrpc interface {
	pb.MerchantServiceServer
}
