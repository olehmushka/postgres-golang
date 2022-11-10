package postgres

import "context"

type Pool interface {
	AcquireConn(ctx context.Context) (PgxConn, error)
	ClosePool()
}
