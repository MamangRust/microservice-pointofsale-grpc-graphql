package handler

import (
	"github.com/MamangRust/microservice-point-of-sale-category/service"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Category pb.CategoryServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Category: NewCategoryHandleGrpc(
			deps.Service,
			deps.Logger,
		),
	}
}
