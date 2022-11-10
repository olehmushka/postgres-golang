package postgres

import (
	"context"

	"github.com/google/uuid"
)

type Writer interface {
	InsertRow(ctx context.Context, query string, args ...interface{}) (uuid.UUID, error)
	InsertRowWithStringID(ctx context.Context, query string, args ...interface{}) (string, error)
	UpdateRowByID(ctx context.Context, sqlData map[string]interface{}, schemaName, tableName, id string) error
	SendBatch(ctx context.Context, batch Batch) error
	ExecuteQuery(ctx context.Context, query string, args ...interface{}) error
	TruncateTables(ctx context.Context, tables []string) error
}
