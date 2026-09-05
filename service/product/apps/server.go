package apps

import (
	db "github.com/MamangRust/microservice-point-of-sale-product/database/schema"
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/server"
	"github.com/MamangRust/microservice-point-of-sale-product/handler"
	mencache "github.com/MamangRust/microservice-point-of-sale-product/cache"
	"github.com/MamangRust/microservice-point-of-sale-product/repository"
	"github.com/MamangRust/microservice-point-of-sale-product/service"
	"github.com/MamangRust/microservice-point-of-sale-shared/observability"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(db.New(srv.Pool))

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	obs := observability.NewTraceLoggerObservability(srv.Logger)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Ctx:           context.Background(),
		Repositories:  repos,
		Logger:        srv.Logger,
		Observability: obs,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterProductServiceServer(gs, handlers.Product)
	}

	return srv, nil
}
