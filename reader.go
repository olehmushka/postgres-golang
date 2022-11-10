package postgres

import "context"

type Reader interface {
	QueryRow(ctx context.Context, query string, args ...any) (Row, error)
	QueryRows(ctx context.Context, conn PgxConn, query string, args ...any) (Rows, error)
	CountRows(ctx context.Context, query string, args ...any) (int, error)
	TruncateTables(ctx context.Context, tables []string) error
}
