package handler

import (
	"context"
	"fmt"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/MamangRust/microservice-point-of-sale-stats-reader/repository"
	"go.uber.org/zap"
)

// OrderStatsHandler serves revenue + order aggregates from ClickHouse
// (order_daily). Results are cached in Redis for 5 minutes.
type OrderStatsHandler struct {
	pb.UnimplementedOrderStatsServiceServer
	repo  repository.Repository
	cache *StatsCache
	log   logger.LoggerInterface
}

func NewOrderStatsHandler(repo repository.Repository, cache *StatsCache, log logger.LoggerInterface) *OrderStatsHandler {
	return &OrderStatsHandler{repo: repo, cache: cache, log: log}
}

func (h *OrderStatsHandler) FindMonthlyTotalRevenue(ctx context.Context, req *pb.FindYearMonthStatsRequest) (*pb.ApiResponseOrderMonthlyRevenue, error) {
	key := fmt.Sprintf("stats:reader:order:monthly-total-revenue:%d:%d", req.GetYear(), req.GetMonth())
	if cached, found := CacheGet[pb.ApiResponseOrderMonthlyRevenue](ctx, h.cache, key); found {
		return cached, nil
	}

	data, err := h.repo.GetMonthlyTotalRevenue(ctx, int(req.GetYear()), int(req.GetMonth()))
	if err != nil {
		h.log.Error("FindMonthlyTotalRevenue failed", zap.Error(err))
		return nil, err
	}

	resp := &pb.ApiResponseOrderMonthlyRevenue{
		Status:  "success",
		Message: "Monthly total revenue retrieved successfully",
		Data:    mapMonthlyRevenue(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *OrderStatsHandler) FindYearlyTotalRevenue(ctx context.Context, req *pb.FindYearStatsRequest) (*pb.ApiResponseOrderYearlyRevenue, error) {
	key := fmt.Sprintf("stats:reader:order:yearly-total-revenue:%d", req.GetYear())
	if cached, found := CacheGet[pb.ApiResponseOrderYearlyRevenue](ctx, h.cache, key); found {
		return cached, nil
	}

	data, err := h.repo.GetYearlyTotalRevenue(ctx, int(req.GetYear()))
	if err != nil {
		h.log.Error("FindYearlyTotalRevenue failed", zap.Error(err))
		return nil, err
	}

	resp := &pb.ApiResponseOrderYearlyRevenue{
		Status:  "success",
		Message: "Yearly total revenue retrieved successfully",
		Data:    mapYearlyRevenue(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *OrderStatsHandler) FindCashierMonthlyRevenue(ctx context.Context, req *pb.FindCashierStatsRequest) (*pb.ApiResponseCashierMonthlyRevenue, error) {
	key := fmt.Sprintf("stats:reader:order:cashier-monthly-revenue:%d", req.GetCashierId())
	if cached, found := CacheGet[pb.ApiResponseCashierMonthlyRevenue](ctx, h.cache, key); found {
		return cached, nil
	}

	data, err := h.repo.GetCashierMonthlyRevenue(ctx, int(req.GetCashierId()))
	if err != nil {
		h.log.Error("FindCashierMonthlyRevenue failed", zap.Error(err))
		return nil, err
	}

	resp := &pb.ApiResponseCashierMonthlyRevenue{
		Status:  "success",
		Message: "Cashier monthly revenue retrieved successfully",
		Data:    mapCashierMonthlyRevenue(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

// --- Mappers ---

func mapMonthlyRevenue(data []repository.MonthlyRevenue) []*pb.OrderMonthlyRevenueResponse {
	var out []*pb.OrderMonthlyRevenueResponse
	for _, d := range data {
		out = append(out, &pb.OrderMonthlyRevenueResponse{
			Year:         d.Year,
			Month:        d.Month,
			TotalRevenue: d.TotalRevenue,
			OrderCount:   int32(d.OrderCount),
		})
	}
	return out
}

func mapYearlyRevenue(data []repository.YearlyRevenue) []*pb.OrderYearlyRevenueResponse {
	var out []*pb.OrderYearlyRevenueResponse
	for _, d := range data {
		out = append(out, &pb.OrderYearlyRevenueResponse{
			Year:         d.Year,
			TotalRevenue: d.TotalRevenue,
			OrderCount:   int32(d.OrderCount),
		})
	}
	return out
}

func mapCashierMonthlyRevenue(data []repository.CashierMonthlyRevenue) []*pb.CashierMonthlyRevenueResponse {
	var out []*pb.CashierMonthlyRevenueResponse
	for _, d := range data {
		out = append(out, &pb.CashierMonthlyRevenueResponse{
			Year:         d.Year,
			Month:        d.Month,
			CashierId:    int32(d.CashierID),
			TotalRevenue: d.TotalRevenue,
			OrderCount:   int32(d.OrderCount),
		})
	}
	return out
}
