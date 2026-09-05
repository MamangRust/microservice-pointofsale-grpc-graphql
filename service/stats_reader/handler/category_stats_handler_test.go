package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/MamangRust/microservice-point-of-sale-stats-reader/repository"
)

var _ logger.LoggerInterface = (*stubLogger)(nil)

func newCategoryHandler(repo repository.Repository) *CategoryStatsHandler {
	return NewCategoryStatsHandler(repo, nil, &stubLogger{})
}

func TestCategoryStatsHandler_FindMonthlySold_Success(t *testing.T) {
	mock := &mockRepo{
		categorySold: []repository.CategoryMonthlySold{
			{Month: "Jan", CategoryID: 1, Quantity: 10, Subtotal: 50000},
			{Month: "Jan", CategoryID: 2, Quantity: 5, Subtotal: 25000},
		},
	}
	h := newCategoryHandler(mock)

	resp, err := h.FindMonthlySold(context.Background(), &pb.FindYearMonthStatsRequest{
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
	if resp.Data[0].CategoryId != 1 {
		t.Errorf("expected category ID 1, got %d", resp.Data[0].CategoryId)
	}
	if resp.Data[0].Quantity != 10 {
		t.Errorf("expected quantity 10, got %d", resp.Data[0].Quantity)
	}
}

func TestCategoryStatsHandler_FindMonthlySold_EmptyResult(t *testing.T) {
	mock := &mockRepo{categorySold: []repository.CategoryMonthlySold{}}
	h := newCategoryHandler(mock)

	resp, err := h.FindMonthlySold(context.Background(), &pb.FindYearMonthStatsRequest{
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

func TestCategoryStatsHandler_FindMonthlySold_Error(t *testing.T) {
	mock := &mockRepo{err: errors.New("clickhouse down")}
	h := newCategoryHandler(mock)

	_, err := h.FindMonthlySold(context.Background(), &pb.FindYearMonthStatsRequest{
		Year:  2026,
		Month: 1,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCategoryStatsHandler_FindMonthlySold_NilData(t *testing.T) {
	mock := &mockRepo{categorySold: nil}
	h := newCategoryHandler(mock)

	resp, err := h.FindMonthlySold(context.Background(), &pb.FindYearMonthStatsRequest{
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
