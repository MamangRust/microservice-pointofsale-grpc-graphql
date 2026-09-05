package handler

import (
	"context"
	"fmt"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/MamangRust/microservice-point-of-sale-stats-reader/repository"
	"go.uber.org/zap"
)

// CashierStatsHandler serves per-cashier order aggregates from ClickHouse
// (order_daily), cached in Redis for 5 minutes.
type CashierStatsHandler struct {
	pb.UnimplementedCashierStatsServiceServer
	repo  repository.Repository
	cache *StatsCache
	log   logger.LoggerInterface
}

func NewCashierStatsHandler(repo repository.Repository, cache *StatsCache, log logger.LoggerInterface) *CashierStatsHandler {
	return &CashierStatsHandler{repo: repo, cache: cache, log: log}
}

func (h *CashierStatsHandler) FindMonthlyOrders(ctx context.Context, req *pb.FindCashierStatsRequest) (*pb.ApiResponseCashierMonthlyOrders, error) {
	key := fmt.Sprintf("stats:reader:cashier:monthly-orders:%d", req.GetCashierId())
	if cached, found := CacheGet[pb.ApiResponseCashierMonthlyOrders](ctx, h.cache, key); found {
		return cached, nil
	}

	data, err := h.repo.GetCashierMonthlyOrders(ctx, int(req.GetCashierId()))
	if err != nil {
		h.log.Error("FindMonthlyOrders failed", zap.Error(err))
		return nil, err
	}

	resp := &pb.ApiResponseCashierMonthlyOrders{
		Status:  "success",
		Message: "Cashier monthly orders retrieved successfully",
		Data:    mapCashierMonthlyOrders(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func mapCashierMonthlyOrders(data []repository.CashierMonthlyOrders) []*pb.CashierMonthlyOrdersResponse {
	var out []*pb.CashierMonthlyOrdersResponse
	for _, d := range data {
		out = append(out, &pb.CashierMonthlyOrdersResponse{
			Month:       d.Month,
			CashierId:   int32(d.CashierID),
			OrderCount:  int64(d.OrderCount),
			TotalAmount: d.TotalAmount,
		})
	}
	return out
}
