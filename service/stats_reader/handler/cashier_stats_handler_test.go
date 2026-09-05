package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/MamangRust/microservice-point-of-sale-stats-reader/repository"
)

func newCashierHandler(repo repository.Repository) *CashierStatsHandler {
	return NewCashierStatsHandler(repo, nil, &stubLogger{})
}

func TestCashierStatsHandler_FindMonthlyOrders_Success(t *testing.T) {
	mock := &mockRepo{
		cashierOrders: []repository.CashierMonthlyOrders{
			{Month: "Jan", CashierID: 1, OrderCount: 30, TotalAmount: 1500000},
			{Month: "Feb", CashierID: 1, OrderCount: 40, TotalAmount: 2000000},
		},
	}
	h := newCashierHandler(mock)

	resp, err := h.FindMonthlyOrders(context.Background(), &pb.FindCashierStatsRequest{
		CashierId: 1,
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
	if resp.Data[0].CashierId != 1 {
		t.Errorf("expected cashier ID 1, got %d", resp.Data[0].CashierId)
	}
	if resp.Data[0].OrderCount != 30 {
		t.Errorf("expected order count 30, got %d", resp.Data[0].OrderCount)
	}
}

func TestCashierStatsHandler_FindMonthlyOrders_Empty(t *testing.T) {
	mock := &mockRepo{cashierOrders: []repository.CashierMonthlyOrders{}}
	h := newCashierHandler(mock)

	resp, err := h.FindMonthlyOrders(context.Background(), &pb.FindCashierStatsRequest{
		CashierId: 99,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 0 {
		t.Errorf("expected 0 data items, got %d", len(resp.Data))
	}
}

func TestCashierStatsHandler_FindMonthlyOrders_Error(t *testing.T) {
	mock := &mockRepo{err: errors.New("clickhouse down")}
	h := newCashierHandler(mock)

	_, err := h.FindMonthlyOrders(context.Background(), &pb.FindCashierStatsRequest{
		CashierId: 1,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCashierStatsHandler_FindMonthlyOrders_NilData(t *testing.T) {
	mock := &mockRepo{cashierOrders: nil}
	h := newCashierHandler(mock)

	resp, err := h.FindMonthlyOrders(context.Background(), &pb.FindCashierStatsRequest{
		CashierId: 1,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 0 {
		t.Errorf("expected 0 data items for nil, got %d", len(resp.Data))
	}
}
