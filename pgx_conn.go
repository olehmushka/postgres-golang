package postgres

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type PgxConn interface {
	// Release returns c to the pool it was acquired from. Once Release has been called, other methods must not be called.
	// However, it is safe to call Release multiple times. Subsequent calls after the first will be ignored.
	Release()

	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)

	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)

	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row

	QueryFunc(ctx context.Context, sql string, args []any, scans []any, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)

	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults

	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)

	Begin(ctx context.Context) (pgx.Tx, error)

	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)

	BeginFunc(ctx context.Context, f func(pgx.Tx) error) error
	Ping(ctx context.Context) error
	Conn() *pgx.Conn
}

type pgxConn interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}
