package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/MamangRust/microservice-point-of-sale-stats-reader/repository"
)

func newTransactionHandler(repo repository.Repository) *TransactionStatsHandler {
	return NewTransactionStatsHandler(repo, nil, &stubLogger{})
}

func TestTransactionStatsHandler_FindMonthlySuccess_Success(t *testing.T) {
	mock := &mockRepo{
		transactionSuccess: []repository.TransactionMonthlySuccess{
			{Month: "Jan", TotalCount: 100, TotalAmount: 5000000},
			{Month: "Feb", TotalCount: 120, TotalAmount: 6000000},
		},
	}
	h := newTransactionHandler(mock)

	resp, err := h.FindMonthlySuccess(context.Background(), &pb.FindYearMonthStatsRequest{
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
	if resp.Data[0].TotalCount != 100 {
		t.Errorf("expected total count 100, got %d", resp.Data[0].TotalCount)
	}
	if resp.Data[0].TotalAmount != 5000000 {
		t.Errorf("expected total amount 5000000, got %d", resp.Data[0].TotalAmount)
	}
}

func TestTransactionStatsHandler_FindMonthlySuccess_Empty(t *testing.T) {
	mock := &mockRepo{transactionSuccess: []repository.TransactionMonthlySuccess{}}
	h := newTransactionHandler(mock)

	resp, err := h.FindMonthlySuccess(context.Background(), &pb.FindYearMonthStatsRequest{
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

func TestTransactionStatsHandler_FindMonthlySuccess_Error(t *testing.T) {
	mock := &mockRepo{err: errors.New("clickhouse down")}
	h := newTransactionHandler(mock)

	_, err := h.FindMonthlySuccess(context.Background(), &pb.FindYearMonthStatsRequest{
		Year:  2026,
		Month: 1,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTransactionStatsHandler_FindMonthlySuccess_NilData(t *testing.T) {
	mock := &mockRepo{transactionSuccess: nil}
	h := newTransactionHandler(mock)

	resp, err := h.FindMonthlySuccess(context.Background(), &pb.FindYearMonthStatsRequest{
		Year:  2026,
		Month: 1,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 0 {
		t.Errorf("expected 0 data items for nil, got %d", len(resp.Data))
	}
}
