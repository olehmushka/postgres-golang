package postgres

import (
	"context"

	"github.com/olehmushka/golang-toolkit/wrapped_error"
)

type transaction struct {
	connection
	poolConn PgxConn
}

func newTransaction(pgxTx pgxTx, poolConn PgxConn) Transaction {
	return &transaction{
		connection: connection{
			pgxConn: pgxTx,
		},
		poolConn: poolConn,
	}
}

func (t *transaction) Commit(ctx context.Context) error {
	defer t.poolConn.Release()
	if err := t.pgxConn.(pgxTx).Commit(ctx); err != nil {
		return wrapped_error.NewInternalServerError(err, "transaction commit was failed")
	}
	return nil
}

func (t *transaction) Rollback(ctx context.Context) error {
	defer t.poolConn.Release()
	if err := t.pgxConn.(pgxTx).Rollback(ctx); err != nil {
		return wrapped_error.NewInternalServerError(err, "transaction rollback was failed")
	}
	return nil
}
