package apps

import (
	"os"

	db "github.com/MamangRust/microservice-point-of-sale-auth/database/schema"

	mencache "github.com/MamangRust/microservice-point-of-sale-auth/cache"
	"github.com/MamangRust/microservice-point-of-sale-auth/handler"
	"github.com/MamangRust/microservice-point-of-sale-auth/repository"
	"github.com/MamangRust/microservice-point-of-sale-auth/service"
	"github.com/MamangRust/microservice-point-of-sale-pkg/auth"
	"github.com/MamangRust/microservice-point-of-sale-pkg/hash"
	"github.com/MamangRust/microservice-point-of-sale-pkg/kafka"
	"github.com/MamangRust/microservice-point-of-sale-pkg/outbox"
	"github.com/MamangRust/microservice-point-of-sale-pkg/server"
	"github.com/MamangRust/microservice-point-of-sale-shared/observability"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	tokenManager, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(db.New(srv.Pool))
	hash := hash.NewHashingPassword()
	// Kafka bersifat opsional: tanpa KAFKA_BROKERS (mis. E2E lokal tanpa kafka)
	// service tetap jalan dan event email di-skip (guard s.kafka != nil).
	var myKafka *kafka.Kafka
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		myKafka = kafka.NewKafka(srv.Logger, []string{brokers})
	}

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	outboxService := outbox.NewOutboxService(srv.Pool, myKafka, srv.Logger)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Repositories:  repos,
		Hash:          hash,
		Token:         tokenManager,
		Logger:        srv.Logger,
		Kafka:         myKafka,
		Pool:          srv.Pool,
		Outbox:        outboxService,
		Observability: observability.NewTraceLoggerObservability(srv.Logger),
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterAuthServiceServer(gs, handlers.Auth)
	}

	// Start the outbox relay so events committed with the business writes are
	// published to Kafka with durable retry and dead-letter semantics.
	go outboxService.Start(srv.Ctx, outbox.OutboxRelayInterval, outbox.OutboxRelayBatchSize)

	return srv, nil
}
