package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/MamangRust/microservice-point-of-sale-shared/domain/events"
)

type mockRepo struct {
	insertOrderErr      error
	insertOrderItemErr  error
	insertTransactionErr error
	closeErr            error
	called              string
}

func (m *mockRepo) InsertOrderEvent(_ context.Context, _ string, _ uint64, _ events.OrderEvent) error {
	m.called = "order"
	return m.insertOrderErr
}

func (m *mockRepo) InsertOrderItemEvent(_ context.Context, _ string, _ uint64, _ events.OrderItemEvent) error {
	m.called = "order_item"
	return m.insertOrderItemErr
}

func (m *mockRepo) InsertTransactionEvent(_ context.Context, _ string, _ uint64, _ events.TransactionEvent) error {
	m.called = "transaction"
	return m.insertTransactionErr
}

func (m *mockRepo) Flush(_ context.Context) error { return nil }
func (m *mockRepo) Close() error                   { return m.closeErr }

func TestUseCase_SaveOrderEvent_Success(t *testing.T) {
	mock := &mockRepo{}
	uc := NewStatsUseCase(mock)

	event := events.OrderEvent{
		OrderID:    1,
		CashierID:  10,
		MerchantID: 20,
		TotalPrice: 50000,
		Status:     "success",
		EventTime:  "2026-01-15T10:00:00Z",
	}

	err := uc.SaveOrderEvent(context.Background(), "evt-001", event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.called != "order" {
		t.Errorf("expected InsertOrderEvent to be called, got %s", mock.called)
	}
}

func TestUseCase_SaveOrderEvent_Error(t *testing.T) {
	mock := &mockRepo{insertOrderErr: errors.New("insert failed")}
	uc := NewStatsUseCase(mock)

	err := uc.SaveOrderEvent(context.Background(), "evt-002", events.OrderEvent{EventTime: "2026-01-15T10:00:00Z"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUseCase_SaveOrderItemEvent_Success(t *testing.T) {
	mock := &mockRepo{}
	uc := NewStatsUseCase(mock)

	event := events.OrderItemEvent{
		OrderItemID: 1,
		OrderID:     10,
		ProductID:   100,
		CategoryID:  5,
		Quantity:    3,
		UnitPrice:   10000,
		Subtotal:    30000,
		EventTime:   "2026-01-15T10:00:00Z",
	}

	err := uc.SaveOrderItemEvent(context.Background(), "evt-003", event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.called != "order_item" {
		t.Errorf("expected InsertOrderItemEvent to be called, got %s", mock.called)
	}
}

func TestUseCase_SaveOrderItemEvent_Error(t *testing.T) {
	mock := &mockRepo{insertOrderItemErr: errors.New("insert failed")}
	uc := NewStatsUseCase(mock)

	err := uc.SaveOrderItemEvent(context.Background(), "evt-004", events.OrderItemEvent{EventTime: "2026-01-15T10:00:00Z"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUseCase_SaveTransactionEvent_Success(t *testing.T) {
	mock := &mockRepo{}
	uc := NewStatsUseCase(mock)

	event := events.TransactionEvent{
		TransactionID: 1,
		OrderID:       10,
		CashierID:     100,
		MerchantID:    200,
		PaymentMethod: "cash",
		Amount:        50000,
		Status:        "success",
		EventTime:     "2026-01-15T10:00:00Z",
	}

	err := uc.SaveTransactionEvent(context.Background(), "evt-005", event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.called != "transaction" {
		t.Errorf("expected InsertTransactionEvent to be called, got %s", mock.called)
	}
}

func TestUseCase_SaveTransactionEvent_Error(t *testing.T) {
	mock := &mockRepo{insertTransactionErr: errors.New("insert failed")}
	uc := NewStatsUseCase(mock)

	err := uc.SaveTransactionEvent(context.Background(), "evt-006", events.TransactionEvent{EventTime: "2026-01-15T10:00:00Z"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUseCase_Close_Success(t *testing.T) {
	mock := &mockRepo{}
	uc := NewStatsUseCase(mock)

	err := uc.Close()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUseCase_Close_Error(t *testing.T) {
	mock := &mockRepo{closeErr: errors.New("close failed")}
	uc := NewStatsUseCase(mock)

	err := uc.Close()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestEventVersion_Empty(t *testing.T) {
	v := eventVersion("")
	if v != 0 {
		t.Errorf("expected 0 for empty string, got %d", v)
	}
}

func TestEventVersion_Valid(t *testing.T) {
	v := eventVersion("2026-01-15T10:00:00Z")
	if v == 0 {
		t.Error("expected non-zero version for valid RFC3339 time")
	}
}

func TestEventVersion_InvalidFormat(t *testing.T) {
	v := eventVersion("not-a-date")
	if v != 0 {
		t.Errorf("expected 0 for invalid format, got %d", v)
	}
}
