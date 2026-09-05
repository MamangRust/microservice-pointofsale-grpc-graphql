package outbox

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInvalidInboxKey = errors.New("invalid consumer inbox key")

// InboxExecutor is the minimal executor surface the consumer-inbox SQL needs.
// Both *pgxpool.Pool and pgx.Tx satisfy it (schema-agnostic, no generated
// schema package dependency).
type InboxExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// ConsumerInbox is the durable deduplication contract used by Kafka handlers
// (Phase 3 — Durable Idempotency). It replaces in-memory-only deduplication:
// reservations survive consumer restarts and rebalances, so at-least-once
// redelivery cannot send the same email twice.
type ConsumerInbox interface {
	Reserve(ctx context.Context, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error)
	MarkProcessed(ctx context.Context, consumerName, eventKey string, reservationVersion int64) error
	Release(ctx context.Context, consumerName, eventKey string, reservationVersion int64, processingErr error) error
}

// PostgresInbox adapts the schema-agnostic consumer-inbox SQL to a live
// *pgxpool.Pool. It keeps the email service's durable-idempotency behavior
// (Phase 3) without depending on a generated schema package.
type PostgresInbox struct {
	pool *pgxpool.Pool
}

func NewPostgresInbox(pool *pgxpool.Pool) (*PostgresInbox, error) {
	if pool == nil {
		return nil, errors.New("consumer inbox requires a non-nil pgx pool")
	}
	return &PostgresInbox{pool: pool}, nil
}

func (p *PostgresInbox) Reserve(ctx context.Context, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error) {
	return Reserve(ctx, p.pool, consumerName, eventKey, topic, partition, offset)
}

func (p *PostgresInbox) MarkProcessed(ctx context.Context, consumerName, eventKey string, reservationVersion int64) error {
	return MarkProcessed(ctx, p.pool, consumerName, eventKey, reservationVersion)
}

func (p *PostgresInbox) Release(ctx context.Context, consumerName, eventKey string, reservationVersion int64, processingErr error) error {
	return Release(ctx, p.pool, consumerName, eventKey, reservationVersion, processingErr)
}

const reserveConsumerInboxSQL = `
WITH reserved AS (
    INSERT INTO consumer_inbox (
        consumer_name, event_key, topic, partition_id, message_offset,
        status, attempts, reservation_version, lease_until, last_error, processed_at
    )
    VALUES ($1, $2, $3, $4, $5, 'processing', 1, 1,
            current_timestamp + interval '1 minute', '', NULL)
    ON CONFLICT (consumer_name, event_key) DO UPDATE
    SET status = 'processing',
        attempts = consumer_inbox.attempts + 1,
        reservation_version = consumer_inbox.reservation_version + 1,
        lease_until = current_timestamp + interval '1 minute',
        last_error = '',
        topic = EXCLUDED.topic,
        partition_id = EXCLUDED.partition_id,
        message_offset = EXCLUDED.message_offset
    WHERE consumer_inbox.status <> 'processed'
      AND consumer_inbox.lease_until <= current_timestamp
    RETURNING reservation_version
)
SELECT
    COALESCE((SELECT reservation_version FROM reserved), 0) AS reservation_version,
    EXISTS (SELECT 1 FROM reserved) AS reserved
`

const consumerInboxStatusSQL = `
SELECT status FROM consumer_inbox WHERE consumer_name = $1 AND event_key = $2
`

const markConsumerInboxProcessedSQL = `
UPDATE consumer_inbox
SET status = 'processed', processed_at = current_timestamp,
    lease_until = current_timestamp, last_error = ''
WHERE consumer_name = $1 AND event_key = $2
  AND status = 'processing' AND reservation_version = $3
`

const releaseConsumerInboxSQL = `
UPDATE consumer_inbox
SET status = 'pending', lease_until = current_timestamp,
    last_error = $3
WHERE consumer_name = $1 AND event_key = $2
  AND status = 'processing' AND reservation_version = $4
`

// Reserve claims an event for a consumer. It returns false when the event was
// already processed. An expired processing lease may be reclaimed after a
// consumer crashes.
func Reserve(ctx context.Context, tx InboxExecutor, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error) {
	if tx == nil || consumerName == "" || eventKey == "" {
		return false, false, 0, ErrInvalidInboxKey
	}
	var reservationVersion int64
	var reserved bool
	err := tx.QueryRow(ctx, reserveConsumerInboxSQL, consumerName, eventKey, topic, partition, offset).
		Scan(&reservationVersion, &reserved)
	if err != nil {
		return false, false, 0, err
	}

	// Determine whether this event was already processed: the upsert only
	// returns a row when it actually (re)reserved the key.
	processed := false
	if !reserved {
		var status string
		if scanErr := tx.QueryRow(ctx, consumerInboxStatusSQL, consumerName, eventKey).Scan(&status); scanErr == nil {
			processed = status == "processed"
		}
	}
	return reserved, processed, reservationVersion, nil
}

func MarkProcessed(ctx context.Context, tx InboxExecutor, consumerName, eventKey string, reservationVersion int64) error {
	if tx == nil || consumerName == "" || eventKey == "" {
		return ErrInvalidInboxKey
	}
	_, err := tx.Exec(ctx, markConsumerInboxProcessedSQL, consumerName, eventKey, reservationVersion)
	return err
}

func Release(ctx context.Context, tx InboxExecutor, consumerName, eventKey string, reservationVersion int64, processingErr error) error {
	if tx == nil || consumerName == "" || eventKey == "" {
		return ErrInvalidInboxKey
	}
	lastError := "consumer processing failed"
	if processingErr != nil {
		lastError = processingErr.Error()
	}
	_, err := tx.Exec(ctx, releaseConsumerInboxSQL, consumerName, eventKey, lastError, reservationVersion)
	return err
}
