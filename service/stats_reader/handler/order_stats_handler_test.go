package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/MamangRust/microservice-point-of-sale-stats-reader/repository"
)

func newOrderHandler(repo repository.Repository) *OrderStatsHandler {
	return NewOrderStatsHandler(repo, nil, &stubLogger{})
}

func TestOrderStatsHandler_FindMonthlyTotalRevenue_Success(t *testing.T) {
	mock := &mockRepo{
		monthlyRevenue: []repository.MonthlyRevenue{
			{Year: "2026", Month: "Jan", TotalRevenue: 1000000, OrderCount: 50},
			{Year: "2026", Month: "Feb", TotalRevenue: 2000000, OrderCount: 80},
		},
	}
	h := newOrderHandler(mock)

	resp, err := h.FindMonthlyTotalRevenue(context.Background(), &pb.FindYearMonthStatsRequest{
		Year:  2026,
		Month: 1,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected status 'success', got '%s'", resp.Status)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 data items, got %d", len(resp.Data))
	}
	if resp.Data[0].TotalRevenue != 1000000 {
		t.Errorf("expected total revenue 1000000, got %d", resp.Data[0].TotalRevenue)
	}
	if resp.Data[0].OrderCount != 50 {
		t.Errorf("expected order count 50, got %d", resp.Data[0].OrderCount)
	}
}

func TestOrderStatsHandler_FindMonthlyTotalRevenue_Empty(t *testing.T) {
	mock := &mockRepo{monthlyRevenue: []repository.MonthlyRevenue{}}
	h := newOrderHandler(mock)

	resp, err := h.FindMonthlyTotalRevenue(context.Background(), &pb.FindYearMonthStatsRequest{
		Year:  2026,
		Month: 1,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 0 {
		t.Errorf("expected 0 data items, got %d", len(resp.Data))
	}
}

func TestOrderStatsHandler_FindMonthlyTotalRevenue_Error(t *testing.T) {
	mock := &mockRepo{err: errors.New("clickhouse down")}
	h := newOrderHandler(mock)

	_, err := h.FindMonthlyTotalRevenue(context.Background(), &pb.FindYearMonthStatsRequest{
		Year:  2026,
		Month: 1,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOrderStatsHandler_FindYearlyTotalRevenue_Success(t *testing.T) {
	mock := &mockRepo{
		yearlyRevenue: []repository.YearlyRevenue{
			{Year: "2026", TotalRevenue: 12000000, OrderCount: 600},
			{Year: "2025", TotalRevenue: 10000000, OrderCount: 500},
		},
	}
	h := newOrderHandler(mock)

	resp, err := h.FindYearlyTotalRevenue(context.Background(), &pb.FindYearStatsRequest{
		Year: 2026,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 data items, got %d", len(resp.Data))
	}
	if resp.Data[0].Year != "2026" {
		t.Errorf("expected year 2026, got '%s'", resp.Data[0].Year)
	}
}

func TestOrderStatsHandler_FindYearlyTotalRevenue_Error(t *testing.T) {
	mock := &mockRepo{err: errors.New("clickhouse down")}
	h := newOrderHandler(mock)

	_, err := h.FindYearlyTotalRevenue(context.Background(), &pb.FindYearStatsRequest{
		Year: 2026,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOrderStatsHandler_FindCashierMonthlyRevenue_Success(t *testing.T) {
	mock := &mockRepo{
		cashierRevenue: []repository.CashierMonthlyRevenue{
			{Year: "2026", Month: "Jan", CashierID: 1, TotalRevenue: 500000, OrderCount: 25},
		},
	}
	h := newOrderHandler(mock)

	resp, err := h.FindCashierMonthlyRevenue(context.Background(), &pb.FindCashierStatsRequest{
		CashierId: 1,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 data item, got %d", len(resp.Data))
	}
	if resp.Data[0].CashierId != 1 {
		t.Errorf("expected cashier ID 1, got %d", resp.Data[0].CashierId)
	}
}

func TestOrderStatsHandler_FindCashierMonthlyRevenue_Empty(t *testing.T) {
	mock := &mockRepo{cashierRevenue: []repository.CashierMonthlyRevenue{}}
	h := newOrderHandler(mock)

	resp, err := h.FindCashierMonthlyRevenue(context.Background(), &pb.FindCashierStatsRequest{
		CashierId: 99,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 0 {
		t.Errorf("expected 0 data items, got %d", len(resp.Data))
	}
}

func TestOrderStatsHandler_FindCashierMonthlyRevenue_Error(t *testing.T) {
	mock := &mockRepo{err: errors.New("clickhouse down")}
	h := newOrderHandler(mock)

	_, err := h.FindCashierMonthlyRevenue(context.Background(), &pb.FindCashierStatsRequest{
		CashierId: 1,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
