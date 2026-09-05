package apps

import (
	db "github.com/MamangRust/microservice-point-of-sale-order/database/schema"
	"context"
	"os"
	"time"

	"github.com/MamangRust/microservice-point-of-sale-order/handler"
	mencache "github.com/MamangRust/microservice-point-of-sale-order/cache"
	"github.com/MamangRust/microservice-point-of-sale-order/repository"
	"github.com/MamangRust/microservice-point-of-sale-order/service"
	"github.com/MamangRust/microservice-point-of-sale-pkg/adapter"
	"github.com/MamangRust/microservice-point-of-sale-pkg/kafka"
	"github.com/MamangRust/microservice-point-of-sale-pkg/outbox"
	"github.com/MamangRust/microservice-point-of-sale-pkg/resilience"
	"github.com/MamangRust/microservice-point-of-sale-pkg/server"
	"github.com/MamangRust/microservice-point-of-sale-shared/cache"
	"github.com/MamangRust/microservice-point-of-sale-shared/observability"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	cashierAddr := getEnv("GRPC_CASHIER_ADDR", "localhost:50055")
	merchantAddr := getEnv("GRPC_MERCHANT_ADDR", "localhost:50056")
	productAddr := getEnv("GRPC_PRODUCT_ADDR", "localhost:50059")
	orderItemAddr := getEnv("GRPC_ORDERITEM_ADDR", "localhost:50057")

	srv.Logger.Info("Connecting to gRPC microservices from Order microservice",
		zap.String("cashier", cashierAddr),
		zap.String("merchant", merchantAddr),
		zap.String("product", productAddr),
		zap.String("order_item", orderItemAddr),
	)

	cashierConn, err := server.NewGRPCClient(cashierAddr)
	if err != nil {
		return nil, err
	}

	merchantConn, err := server.NewGRPCClient(merchantAddr)
	if err != nil {
		cashierConn.Close()
		return nil, err
	}

	productConn, err := server.NewGRPCClient(productAddr)
	if err != nil {
		cashierConn.Close()
		merchantConn.Close()
		return nil, err
	}

	orderItemConn, err := server.NewGRPCClient(orderItemAddr)
	if err != nil {
		cashierConn.Close()
		merchantConn.Close()
		productConn.Close()
		return nil, err
	}

	go func() {
		<-srv.Ctx.Done()
		srv.Logger.Info("Closing gRPC client connections in Order microservice")
		cashierConn.Close()
		merchantConn.Close()
		productConn.Close()
		orderItemConn.Close()
	}()

	cashierClient := pb.NewCashierServiceClient(cashierConn)
	merchantClient := pb.NewMerchantServiceClient(merchantConn)
	productClient := pb.NewProductServiceClient(productConn)
	orderItemClient := pb.NewOrderItemServiceClient(orderItemConn)

	// F6 §8.1 poin 5 — dependency guard: per-call deadline + circuit breaker +
	// bulkhead untuk setiap dependency gRPC (cashier/merchant/product/order_item).
	// Threshold kecil supaya dependency yang down cepat di-fail-fast, bukan di-hammer.
	guardCashier := resilience.NewDependencyGuard("cashier", 5, 30, 100, 3*time.Second, srv.Logger)
	guardMerchant := resilience.NewDependencyGuard("merchant", 5, 30, 100, 3*time.Second, srv.Logger)
	guardProduct := resilience.NewDependencyGuard("product", 5, 30, 100, 3*time.Second, srv.Logger)
	guardOrderItem := resilience.NewDependencyGuard("order_item", 5, 30, 100, 3*time.Second, srv.Logger)

	repos := repository.NewRepositories(
		db.New(srv.Pool), cashierClient, merchantClient, productClient, orderItemClient,
		repository.GuardOptions{
			Cashier:   []adapter.GuardOption{adapter.WithDependencyGuard(guardCashier)},
			Merchant:  []adapter.GuardOption{adapter.WithDependencyGuard(guardMerchant)},
			Product:   []adapter.GuardOption{adapter.WithDependencyGuard(guardProduct)},
			OrderItem: []adapter.GuardOption{adapter.WithDependencyGuard(guardOrderItem)},
		},
	)
	// Kafka bersifat opsional: tanpa KAFKA_BROKERS (mis. E2E lokal tanpa kafka)
	// service tetap jalan dan event stats di-skip (guard s.outbox != nil).
	var myKafka *kafka.Kafka
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		myKafka = kafka.NewKafka(srv.Logger, []string{brokers})
	}
	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	// F3 — transactional outbox: order/order_item stats events are enqueued
	// after the cross-service order flow commits and the relay publishes them
	// durably to the stats.pos.*.event topics (best-effort, like ecommerce).
	outboxService := outbox.NewOutboxService(srv.Pool, myKafka, srv.Logger)

	cacheMetrics, err := observability.NewCacheMetrics("order")
	if err != nil {
		return nil, err
	}

	cacheStore := cache.NewCacheStore(srv.Redis, srv.Logger, cacheMetrics)
	mencacheObj := mencache.NewMencache(cacheStore)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Ctx:           context.Background(),
		Repositories:  repos,
		Pool:          srv.Pool,
		Outbox:        outboxService,
		Logger:        srv.Logger,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterOrderServiceServer(gs, handlers.Order)
	}

	// Start the outbox relay so stats events committed with the business flow
	// are published to Kafka with durable retry and dead-letter semantics.
	go outboxService.Start(srv.Ctx, outbox.OutboxRelayInterval, outbox.OutboxRelayBatchSize)

	return srv, nil
}
