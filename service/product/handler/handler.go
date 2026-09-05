package handler

import (
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-product/service"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Product pb.ProductServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Product: NewProductHandleGrpc(
			deps.Service,
			deps.Logger,
		),
	}
}
