package postgres

import "context"

type ClientWriter interface {
	Writer
	Pool
	BeginTx(ctx context.Context, opts TxOptions) (Transaction, error)
}
