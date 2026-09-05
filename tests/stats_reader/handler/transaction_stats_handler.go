package handler

import (
	"context"
	"fmt"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/MamangRust/microservice-point-of-sale-stats-reader/repository"
	"go.uber.org/zap"
)

// TransactionStatsHandler serves transaction aggregates from ClickHouse
// (transaction_daily), cached in Redis for 5 minutes.
type TransactionStatsHandler struct {
	pb.UnimplementedTransactionStatsServiceServer
	repo  repository.Repository
	cache *StatsCache
	log   logger.LoggerInterface
}

func NewTransactionStatsHandler(repo repository.Repository, cache *StatsCache, log logger.LoggerInterface) *TransactionStatsHandler {
	return &TransactionStatsHandler{repo: repo, cache: cache, log: log}
}

func (h *TransactionStatsHandler) FindMonthlySuccess(ctx context.Context, req *pb.FindYearMonthStatsRequest) (*pb.ApiResponseTransactionMonthlySuccess, error) {
	key := fmt.Sprintf("stats:reader:transaction:monthly-success:%d:%d", req.GetYear(), req.GetMonth())
	if cached, found := CacheGet[pb.ApiResponseTransactionMonthlySuccess](ctx, h.cache, key); found {
		return cached, nil
	}

	data, err := h.repo.GetTransactionMonthlySuccess(ctx, int(req.GetYear()), int(req.GetMonth()))
	if err != nil {
		h.log.Error("FindMonthlySuccess failed", zap.Error(err))
		return nil, err
	}

	resp := &pb.ApiResponseTransactionMonthlySuccess{
		Status:  "success",
		Message: "Monthly transaction success retrieved successfully",
		Data:    mapTransactionMonthlySuccess(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func mapTransactionMonthlySuccess(data []repository.TransactionMonthlySuccess) []*pb.TransactionMonthlySuccessResponse {
	var out []*pb.TransactionMonthlySuccessResponse
	for _, d := range data {
		out = append(out, &pb.TransactionMonthlySuccessResponse{
			Month:       d.Month,
			TotalCount:  int64(d.TotalCount),
			TotalAmount: d.TotalAmount,
		})
	}
	return out
}
