package handler

import (
	"github.com/MamangRust/microservice-point-of-sale-order/service"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Order pb.OrderServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Order: NewOrderHandleGrpc(
			deps.Service,
			deps.Logger,
		),
	}
}
