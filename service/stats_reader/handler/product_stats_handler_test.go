package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/MamangRust/microservice-point-of-sale-stats-reader/repository"
)

func newProductHandler(repo repository.Repository) *ProductStatsHandler {
	return NewProductStatsHandler(repo, nil, &stubLogger{})
}

func TestProductStatsHandler_FindMonthlySold_Success(t *testing.T) {
	mock := &mockRepo{
		productSold: []repository.ProductMonthlySold{
			{Month: "Jan", ProductID: 1, Quantity: 20, Subtotal: 800000},
			{Month: "Jan", ProductID: 2, Quantity: 15, Subtotal: 600000},
		},
	}
	h := newProductHandler(mock)

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
	if resp.Data[0].ProductId != 1 {
		t.Errorf("expected product ID 1, got %d", resp.Data[0].ProductId)
	}
	if resp.Data[0].Quantity != 20 {
		t.Errorf("expected quantity 20, got %d", resp.Data[0].Quantity)
	}
}

func TestProductStatsHandler_FindMonthlySold_Empty(t *testing.T) {
	mock := &mockRepo{productSold: []repository.ProductMonthlySold{}}
	h := newProductHandler(mock)

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

func TestProductStatsHandler_FindMonthlySold_Error(t *testing.T) {
	mock := &mockRepo{err: errors.New("clickhouse down")}
	h := newProductHandler(mock)

	_, err := h.FindMonthlySold(context.Background(), &pb.FindYearMonthStatsRequest{
		Year:  2026,
		Month: 1,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestProductStatsHandler_FindMonthlySold_NilData(t *testing.T) {
	mock := &mockRepo{productSold: nil}
	h := newProductHandler(mock)

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
