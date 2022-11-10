package postgres

import "context"

type pgxTx interface {
	pgxConn
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
