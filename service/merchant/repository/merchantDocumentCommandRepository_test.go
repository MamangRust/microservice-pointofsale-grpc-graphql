package repository

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	db "github.com/MamangRust/microservice-point-of-sale-merchant/database/schema"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-point-of-sale-shared/errors"
)

// fakeDBTX is a hand-rolled db.DBTX that records the executed SQL and its
// positional args and returns a canned row for QueryRow. It allows testing
// the repository param building without a real database.
type fakeDBTX struct {
	lastSQL  string
	lastArgs []any

	row    pgx.Row
	rowErr error
}

func (f *fakeDBTX) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	f.lastSQL = sql
	f.lastArgs = args
	return pgconn.CommandTag{}, nil
}

func (f *fakeDBTX) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	f.lastSQL = sql
	f.lastArgs = args
	return nil, nil
}

func (f *fakeDBTX) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	f.lastSQL = sql
	f.lastArgs = args
	if f.rowErr != nil {
		return fakeRow{err: f.rowErr}
	}
	return f.row
}

// fakeRow is a minimal pgx.Row that returns canned values.
type fakeRow struct {
	vals []any
	err  error
}

func (r fakeRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(r.vals) != len(dest) {
		return fmt.Errorf("fakeRow: got %d dest, want %d values", len(dest), len(r.vals))
	}
	for i, d := range dest {
		if err := setScanDest(d, r.vals[i]); err != nil {
			return err
		}
	}
	return nil
}

// setScanDest copies val into the typed destination d. It covers the concrete
// types produced by the generated MerchantDocument scan (int32, string,
// *string and pgtype.Timestamp).
func setScanDest(dest, val any) error {
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("fakeRow: dest %T is not a non-nil pointer", dest)
	}
	ev := rv.Elem()

	switch v := val.(type) {
	case int32:
		ev.SetInt(int64(v))
	case string:
		switch ev.Interface().(type) {
		case string:
			ev.SetString(v)
		case *string:
			ev.Set(reflect.ValueOf(&v))
		default:
			return fmt.Errorf("fakeRow: unsupported string dest %T", dest)
		}
	case *string:
		if ev.Type() != reflect.TypeOf(v) {
			return fmt.Errorf("fakeRow: cannot set %T into %s", v, ev.Type())
		}
		ev.Set(reflect.ValueOf(v))
	case pgtype.Timestamp:
		ev.Set(reflect.ValueOf(v))
	default:
		return fmt.Errorf("fakeRow: unsupported value type %T", val)
	}
	return nil
}

// merchantDocumentRowFixture mirrors the column order of the generated
// MerchantDocument scan (RETURNING *).
func merchantDocumentRowFixture(docID, merchantID int) []any {
	note := "Approved"
	return []any{
		int32(docID),
		int32(merchantID),
		"updated_type",
		"https://example.com/updated.pdf",
		"verified",
		&note,
		pgtype.Timestamp{},
		pgtype.Timestamp{},
		pgtype.Timestamp{},
		pgtype.Timestamp{},
	}
}

// Regression tests for the DocumentID-vs-MerchantID fix: the update must be
// keyed on the request's DocumentID — never on the MerchantID. The fixture
// deliberately uses different values (DocumentID=42, MerchantID=999) so a
// regression (using MerchantID as the WHERE key) cannot pass by coincidence.

func TestUpdateMerchantDocument_UsesDocumentID_NotMerchantID(t *testing.T) {
	const (
		docID      = 42
		merchantID = 999
	)

	fake := &fakeDBTX{row: fakeRow{vals: merchantDocumentRowFixture(docID, merchantID)}}
	repo := NewMerchantDocumentCommandRepository(db.New(fake))

	updated, err := repo.UpdateMerchantDocument(context.Background(), &requests.UpdateMerchantDocumentRequest{
		DocumentID:   ptr(docID),
		MerchantID:   merchantID,
		DocumentType: "updated_type",
		DocumentUrl:  "https://example.com/updated.pdf",
		Status:       "verified",
		Note:         "Approved",
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, int32(docID), updated.DocumentID)
	assert.Equal(t, int32(merchantID), updated.MerchantID)

	// The first positional arg ($1 — the WHERE key) must be the DocumentID.
	require.Len(t, fake.lastArgs, 5)
	assert.Equal(t, int32(docID), fake.lastArgs[0],
		"WHERE key must be the request DocumentID, never the MerchantID")
	assert.NotEqual(t, int32(merchantID), fake.lastArgs[0])

	// All params map to the generated query in order: DocumentID, DocumentType,
	// DocumentUrl, Status, Note.
	assert.Equal(t, "updated_type", fake.lastArgs[1])
	assert.Equal(t, "https://example.com/updated.pdf", fake.lastArgs[2])
	assert.Equal(t, "verified", fake.lastArgs[3])
	noteArg, ok := fake.lastArgs[4].(*string)
	require.True(t, ok)
	assert.Equal(t, "Approved", *noteArg)

	// The generated SQL must filter by document_id — a "merchant_id = $n"
	// predicate is exactly the bug this fix removes.
	assert.Contains(t, fake.lastSQL, "document_id = $1 AND deleted_at IS NULL")
	assert.NotContains(t, fake.lastSQL, "merchant_id = $")
}

func TestUpdateMerchantDocumentStatus_UsesDocumentID_NotMerchantID(t *testing.T) {
	const (
		docID      = 42
		merchantID = 999
	)

	fake := &fakeDBTX{row: fakeRow{vals: merchantDocumentRowFixture(docID, merchantID)}}
	repo := NewMerchantDocumentCommandRepository(db.New(fake))

	updated, err := repo.UpdateMerchantDocumentStatus(context.Background(), &requests.UpdateMerchantDocumentStatusRequest{
		DocumentID: ptr(docID),
		MerchantID: merchantID,
		Status:     "rejected",
		Note:       "Invalid document",
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, int32(docID), updated.DocumentID)

	require.Len(t, fake.lastArgs, 3)
	assert.Equal(t, int32(docID), fake.lastArgs[0],
		"WHERE key must be the request DocumentID, never the MerchantID")
	assert.NotEqual(t, int32(merchantID), fake.lastArgs[0])

	// Remaining params map in order: Status, Note.
	assert.Equal(t, "rejected", fake.lastArgs[1])
	noteArg, ok := fake.lastArgs[2].(*string)
	require.True(t, ok)
	assert.Equal(t, "Invalid document", *noteArg)

	assert.Contains(t, fake.lastSQL, "document_id = $1 AND deleted_at IS NULL")
	assert.NotContains(t, fake.lastSQL, "merchant_id = $")
}

func TestUpdateMerchantDocument_NilDocumentID_ReturnsBadRequest(t *testing.T) {
	fake := &fakeDBTX{}
	repo := NewMerchantDocumentCommandRepository(db.New(fake))

	_, err := repo.UpdateMerchantDocument(context.Background(), &requests.UpdateMerchantDocumentRequest{
		MerchantID: 999,
	})
	require.Error(t, err)

	var appErr *sharedErrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, sharedErrors.ErrorTypeBadRequest, appErr.Type)
	assert.Contains(t, appErr.Message, "document id is required")

	// The DB must never be touched when DocumentID is missing.
	assert.Empty(t, fake.lastSQL)
}

func TestUpdateMerchantDocumentStatus_NilDocumentID_ReturnsBadRequest(t *testing.T) {
	fake := &fakeDBTX{}
	repo := NewMerchantDocumentCommandRepository(db.New(fake))

	_, err := repo.UpdateMerchantDocumentStatus(context.Background(), &requests.UpdateMerchantDocumentStatusRequest{
		MerchantID: 999,
	})
	require.Error(t, err)

	var appErr *sharedErrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, sharedErrors.ErrorTypeBadRequest, appErr.Type)
	assert.Contains(t, appErr.Message, "document id is required")

	assert.Empty(t, fake.lastSQL)
}

func TestUpdateMerchantDocument_DBError_WrapsInternal(t *testing.T) {
	inner := errors.New("connection refused")
	fake := &fakeDBTX{rowErr: inner}
	repo := NewMerchantDocumentCommandRepository(db.New(fake))

	_, err := repo.UpdateMerchantDocument(context.Background(), &requests.UpdateMerchantDocumentRequest{
		DocumentID: ptr(42),
		MerchantID: 999,
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, sharedErrors.ErrInternal)
	assert.ErrorIs(t, err, inner)
}

func TestUpdateMerchantDocumentStatus_DBError_WrapsInternal(t *testing.T) {
	inner := errors.New("connection refused")
	fake := &fakeDBTX{rowErr: inner}
	repo := NewMerchantDocumentCommandRepository(db.New(fake))

	_, err := repo.UpdateMerchantDocumentStatus(context.Background(), &requests.UpdateMerchantDocumentStatusRequest{
		DocumentID: ptr(42),
		MerchantID: 999,
		Status:     "approved",
		Note:       "OK",
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, sharedErrors.ErrInternal)
	assert.ErrorIs(t, err, inner)
}
