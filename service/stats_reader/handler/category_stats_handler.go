package handler

import (
	"context"
	"fmt"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/MamangRust/microservice-point-of-sale-stats-reader/repository"
	"go.uber.org/zap"
)

// CategoryStatsHandler serves per-category sold aggregates from ClickHouse
// (order_item_daily), cached in Redis for 5 minutes.
type CategoryStatsHandler struct {
	pb.UnimplementedCategoryStatsServiceServer
	repo  repository.Repository
	cache *StatsCache
	log   logger.LoggerInterface
}

func NewCategoryStatsHandler(repo repository.Repository, cache *StatsCache, log logger.LoggerInterface) *CategoryStatsHandler {
	return &CategoryStatsHandler{repo: repo, cache: cache, log: log}
}

func (h *CategoryStatsHandler) FindMonthlySold(ctx context.Context, req *pb.FindYearMonthStatsRequest) (*pb.ApiResponseCategoryMonthlySold, error) {
	key := fmt.Sprintf("stats:reader:category:monthly-sold:%d:%d", req.GetYear(), req.GetMonth())
	if cached, found := CacheGet[pb.ApiResponseCategoryMonthlySold](ctx, h.cache, key); found {
		return cached, nil
	}

	data, err := h.repo.GetCategoryMonthlySold(ctx, int(req.GetYear()), int(req.GetMonth()))
	if err != nil {
		h.log.Error("FindMonthlySold failed", zap.Error(err))
		return nil, err
	}

	resp := &pb.ApiResponseCategoryMonthlySold{
		Status:  "success",
		Message: "Monthly category sold retrieved successfully",
		Data:    mapCategoryMonthlySold(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func mapCategoryMonthlySold(data []repository.CategoryMonthlySold) []*pb.CategoryMonthlySoldResponse {
	var out []*pb.CategoryMonthlySoldResponse
	for _, d := range data {
		out = append(out, &pb.CategoryMonthlySoldResponse{
			Month:      d.Month,
			CategoryId: int32(d.CategoryID),
			Quantity:   int64(d.Quantity),
			Subtotal:   d.Subtotal,
		})
	}
	return out
}
