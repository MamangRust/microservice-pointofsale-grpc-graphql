package handler

import (
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-role/service"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Role pb.RoleServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Role: NewRoleHandleGrpc(
			deps.Service,
			deps.Logger,
		),
	}
}
