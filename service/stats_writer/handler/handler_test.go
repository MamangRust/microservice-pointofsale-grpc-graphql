package handler

import (
	"encoding/json"
	"testing"
	"time"
)

func TestStatsTopics(t *testing.T) {
	topics := StatsTopics()
	if len(topics) != 3 {
		t.Fatalf("expected 3 topics, got %d", len(topics))
	}
	expected := map[string]bool{
		"stats.pos.order.event":       false,
		"stats.pos.order_item.event":  false,
		"stats.pos.transaction.event": false,
	}
	for _, topic := range topics {
		if _, ok := expected[topic]; !ok {
			t.Errorf("unexpected topic: %s", topic)
		}
		expected[topic] = true
	}
	for topic, found := range expected {
		if !found {
			t.Errorf("missing topic: %s", topic)
		}
	}
}

func TestStatsHandler_TryUnwrap_ValidEnvelope(t *testing.T) {
	env := statEnvelope{
		EventID: "evt-123",
		Payload: json.RawMessage(`{"order_id": 1}`),
	}
	raw, _ := json.Marshal(env)

	h := &StatsHandler{dedup: make(map[string]time.Time)}
	result := h.tryUnwrap(raw)

	if result == nil {
		t.Fatal("expected non-nil envelope")
	}
	if result.EventID != "evt-123" {
		t.Errorf("expected event ID 'evt-123', got '%s'", result.EventID)
	}
}

func TestStatsHandler_TryUnwrap_MissingEventID(t *testing.T) {
	env := statEnvelope{
		EventID: "",
		Payload: json.RawMessage(`{"order_id": 1}`),
	}
	raw, _ := json.Marshal(env)

	h := &StatsHandler{dedup: make(map[string]time.Time)}
	result := h.tryUnwrap(raw)

	if result != nil {
		t.Fatal("expected nil for missing event ID")
	}
}

func TestStatsHandler_TryUnwrap_InvalidJSON(t *testing.T) {
	h := &StatsHandler{dedup: make(map[string]time.Time)}
	result := h.tryUnwrap([]byte("not json"))

	if result != nil {
		t.Fatal("expected nil for invalid JSON")
	}
}

func TestStatsHandler_TryUnwrap_RawPayload(t *testing.T) {
	// A plain payload without envelope wrapper
	raw := []byte(`{"order_id": 1}`)

	h := &StatsHandler{dedup: make(map[string]time.Time)}
	result := h.tryUnwrap(raw)

	if result != nil {
		t.Fatal("expected nil for raw payload without event_id")
	}
}

func TestStatsHandler_IsDuplicate_FirstSeen(t *testing.T) {
	h := &StatsHandler{dedup: make(map[string]time.Time)}

	if h.isDuplicate("evt-001") {
		t.Error("expected false for first-time event")
	}
}

func TestStatsHandler_IsDuplicate_DuplicateSeen(t *testing.T) {
	h := &StatsHandler{dedup: make(map[string]time.Time)}

	h.isDuplicate("evt-002")
	if !h.isDuplicate("evt-002") {
		t.Error("expected true for duplicate event")
	}
}

func TestStatsHandler_IsDuplicate_EmptyEventID(t *testing.T) {
	h := &StatsHandler{dedup: make(map[string]time.Time)}

	if h.isDuplicate("") {
		t.Error("expected false for empty event ID")
	}
}

func TestStatsHandler_IsDuplicate_OldEventExpired(t *testing.T) {
	h := &StatsHandler{dedup: make(map[string]time.Time)}

	// Insert an event that is older than 24h
	h.dedup["old-event"] = time.Now().Add(-25 * time.Hour)

	// Should NOT be considered duplicate (it expired), and the event is
	// re-added with current timestamp (isDuplicate records it after sweep).
	if h.isDuplicate("old-event") {
		t.Error("expected false for expired event")
	}
	// Verify the old entry was cleaned and replaced with current time
	if _, ok := h.dedup["old-event"]; !ok {
		t.Error("expected old event to be re-recorded with current timestamp")
	}
}
