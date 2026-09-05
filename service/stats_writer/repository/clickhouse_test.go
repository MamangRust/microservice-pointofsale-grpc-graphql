package repository

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestToUUID_ValidUUID(t *testing.T) {
	input := "550e8400-e29b-41d4-a716-446655440000"
	result := toUUID(input)

	expected := uuid.Must(uuid.Parse(input))
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestToUUID_EmptyString(t *testing.T) {
	result := toUUID("")
	if result != uuid.Nil {
		t.Errorf("expected nil UUID, got %v", result)
	}
}

func TestToUUID_InvalidFormat(t *testing.T) {
	result := toUUID("not-a-uuid")
	if result != uuid.Nil {
		t.Errorf("expected nil UUID for invalid format, got %v", result)
	}
}

func TestParseEventTime_ValidRFC3339(t *testing.T) {
	input := "2026-01-15T10:30:00Z"
	result := parseEventTime(input)

	expected := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestParseEventTime_EmptyString(t *testing.T) {
	before := time.Now()
	result := parseEventTime("")
	after := time.Now()

	if result.Before(before) || result.After(after) {
		t.Errorf("expected time near now, got %v", result)
	}
}

func TestParseEventTime_InvalidFormat(t *testing.T) {
	before := time.Now()
	result := parseEventTime("not-a-date")
	after := time.Now()

	if result.Before(before.Add(-time.Second)) || result.After(after.Add(time.Second)) {
		t.Errorf("expected fallback to now for invalid format, got %v", result)
	}
}
